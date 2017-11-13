package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/moolen/asdf/changelog"
	"github.com/moolen/asdf/config"
	"github.com/moolen/asdf/fetcher"
	"github.com/moolen/asdf/repository"
)

// ErrNoCommits is returned if there are no changes between the last release and the current HEAD
var ErrNoCommits = errors.New("there is nothing to release; no new commits found")

// ReleaseToken is repleaced with the prerelease number
// If there was no previous release it will starting with 1
var ReleaseToken = "{RELEASE_NUMBER}"

// CommitToken is replaced within a release and contains the short commit hash
var CommitToken = "{COMMIT_SHA}"

// pullRequestTitleRegex is used to strip a Ticket ID from the PullRequest title
var pullRequestTitleRegex = regexp.MustCompile("(\\w*-[0-9]+)")

func generateRelease(cwd, token, branch string, config *config.Config) error {
	log.Printf("generating release in dir: %s", cwd)
	versionPath := path.Join(cwd, config.VersionFile)
	changelogfile := path.Join(cwd, config.ChangelogFile)
	execDir(cwd, "git", "fetch", "--all")
	fetcher, err := fetcher.New(token, config.Repository)
	if err != nil {
		return err
	}
	changelog, nextVersion, err := generateReleaseAndChangelog(cwd, branch, fetcher, config)
	if err != nil {
		return err
	}
	currentChangelog, err := ioutil.ReadFile(changelogfile)
	_, ok := err.(*os.PathError)
	if err != nil && !ok {
		return err
	}
	if nextVersion == nil {
		return errors.New("could not calculate next version")
	}
	err = ioutil.WriteFile(changelogfile, []byte(fmt.Sprintf("%s\n\n\n%s", changelog, currentChangelog)), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(versionPath, []byte(nextVersion.String()), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func generateReleaseAndChangelog(cwd, branch string, fetcher fetcher.PullRequestFetcher, config *config.Config) (string, *semver.Version, error) {
	log.Println("generate release..")
	var formatter changelog.FormatFunc
	versionPath := path.Join(cwd, config.VersionFile)
	versionFile, err := os.Open(versionPath)
	defer versionFile.Close()
	if err != nil {
		return "", nil, errors.New("version file does not exist, please create one")
	}
	version, err := readVersion(versionFile)
	if err != nil {
		return "", nil, errors.New("version file does not contain a semver version")
	}
	log.Printf("found version in file: %s", version)
	repo := repository.New(cwd, repository.DefaultMapFunc)
	latestReleaseCommit, err := repo.LatestChangeOfFile(path.Base(versionPath))
	if err != nil {
		return "", nil, err
	}
	log.Printf("latest release commit: (%s) %s", latestReleaseCommit.Hash, latestReleaseCommit.Subject)
	commits, err := repo.GetHistoryUntil(latestReleaseCommit)
	if err != nil {
		return "", nil, err
	}
	if len(commits) == 0 {
		return "", nil, ErrNoCommits
	}
	log.Printf("found %d commits since last release commit", len(commits))
	nextVersion, err := calcNextVersion(commits[0], branch, version, config.BranchSuffix, commits.MaxChange())
	if err != nil {
		return "", nil, err
	}
	log.Printf("next version: %s", nextVersion)

	if fetcher != nil {
		formatter, err = createPRFormatter(fetcher, config.TicketURL)
		if err != nil {
			return "", nil, err
		}
	} else {
		formatter = changelog.DefaultFormatFunc
	}

	cl := changelog.New(config.Types, formatter)
	changelog := cl.Create(commits, nextVersion)
	return changelog, nextVersion, nil
}

// this returns a FormatFunc for the changelog.
// it will check if PRs matches the Ticket ID from the commit
// and renders the PR ID in the changelog
func createPRFormatter(fetcher fetcher.PullRequestFetcher, url string) (changelog.FormatFunc, error) {
	pullRequests, err := fetcher.Fetch()
	PullRequestMap := make(map[string][]string)
	if err != nil {
		return nil, err
	}
	formatPullRequestID := func(ID int) string {
		return fmt.Sprintf("#%d", ID)
	}
	// parse all pull requests and put them into a map
	// so we can have a easy direct lookup
	for _, pr := range pullRequests {
		if pr.Merged == true && pr.Title != "" {
			matched := pullRequestTitleRegex.FindAllStringSubmatch(pr.Title, 1)
			if len(matched) > 0 {
				ticketID := matched[0][0]
				if PullRequestMap[ticketID] == nil {
					PullRequestMap[ticketID] = []string{formatPullRequestID(pr.ID)}
				} else {
					PullRequestMap[ticketID] = append(PullRequestMap[ticketID], formatPullRequestID(pr.ID))
				}
			}
		}
	}
	// return the changelog.FormatFunc
	return func(c *repository.Commit) string {
		if c.Scope != "" {
			var prList string
			ticketURL := strings.Replace(url, "{SCOPE}", c.Scope, -1)
			if len(PullRequestMap[c.Scope]) > 0 {
				prList = strings.Join(PullRequestMap[c.Scope], ", ")
				return fmt.Sprintf("* %s [%s](%s) (%s) \n", c.Subject, c.Scope, ticketURL, prList)
			}
			return fmt.Sprintf("* %s [%s](%s) \n", c.Subject, c.Scope, ticketURL)
		}
		return changelog.DefaultFormatFunc(c)
	}, nil
}

func calcNextVersion(latestCommit *repository.Commit, branch string, latest *semver.Version, branchSuffix map[string]string, change repository.Change) (*semver.Version, error) {
	var next semver.Version
	var err error
	log.Printf("latest release: %s", latest)
	// if we're on master we'll either:
	// - remove the prerelease suffix
	// - or increment to next version
	if branch == "master" {
		if latest.Prerelease() == "" {
			log.Printf("on branch master without prerelease, calculating next change: %d", change)
			next = nextReleaseByChange(latest, change)
		} else {
			next, err = latest.SetPrerelease("")
			log.Printf("on branch master with prerelease, removing prerelease")
			if err != nil {
				return nil, err
			}
		}
		return &next, nil
	}

	// if this branch has a special mapping
	// we apply the mapping withouth incrementing
	for branchRE, suffix := range branchSuffix {
		re := regexp.MustCompile(branchRE)
		if re.MatchString(branch) {
			nextSuffix := nextPrereleaseSuffix(latest, latestCommit, suffix)
			next, err = latest.SetPrerelease(nextSuffix)
		}
	}
	return &next, err
}

func nextPrereleaseSuffix(latest *semver.Version, commit *repository.Commit, suffix string) string {
	var nextSuffix string
	hash := commit.Hash
	if len(hash) > 8 {
		hash = hash[:8]
	}
	prerelease := latest.Prerelease()
	if strings.Contains(suffix, ReleaseToken) {
		stripped := strings.Replace(suffix, ReleaseToken, "", 1)
		strippedPrerelease := strings.Replace(prerelease, stripped, "", 1)
		releaseNumber, _ := strconv.ParseInt(strippedPrerelease, 10, 64)
		nextReleaseNum := fmt.Sprintf("%d", releaseNumber+1)
		nextSuffix = strings.Replace(suffix, ReleaseToken, nextReleaseNum, 1)
	} else if strings.Contains(suffix, CommitToken) {
		nextSuffix = strings.Replace(suffix, "{COMMIT_SHA}", hash, 1)
	}
	return nextSuffix
}

func nextReleaseByChange(latest *semver.Version, change repository.Change) semver.Version {
	switch change {
	case repository.MajorChange:
		log.Printf("increment major")
		return latest.IncMajor()
	case repository.MinorChange:
		log.Printf("increment minor")
		return latest.IncMinor()
	}
	return latest.IncPatch()
}

func execDir(dir, cmd string, things ...string) {
	c := exec.Command(cmd, things...)
	c.Dir = dir
	err := c.Run()
	if err != nil {
		fmt.Printf("dir: %s, cmd: %s, args: %s", dir, cmd, things)
		panic(err)
	}
}
