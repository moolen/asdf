package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/moolen/asdf/changelog"
	"github.com/moolen/asdf/config"
	"github.com/moolen/asdf/repository"
	"github.com/urfave/cli"
)

// ErrNoCommits is returned if there are no changes between the last release and the current HEAD
var ErrNoCommits = errors.New("there is nothing to release: no new commits found")

// ErrNoRevision the use has to specify a revision(range)
// set man 7 gitrevisions
var ErrNoRevision = errors.New("revision is required")

// ReleaseToken is repleaced with the prerelease number
// If there was no previous release it will starting with 1
var ReleaseToken = "{RELEASE_NUMBER}"

// CommitToken is replaced within a release and contains the short commit hash
var CommitToken = "{COMMIT_SHA}"

// pullRequestTitleRegex is used to strip a Ticket ID from the PullRequest title
var pullRequestTitleRegex = regexp.MustCompile("(\\w*-[0-9]+)")

// generateCommand is a stateful operation
// It looks for a VERSION file, calculates the
// changelog based on the commits since this file has changed
func generateCommand(c *cli.Context) error {
	branch := c.GlobalString(flagBranch)
	token := c.GlobalString(flagGithubToken)
	cwd, err := getCwd(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	config, err := config.FromFile(path.Join(cwd, configFilename))
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	log.Printf("generating release in dir: %s", cwd)
	versionPath := path.Join(cwd, config.VersionFile)
	changelogfile := path.Join(cwd, config.ChangelogFile)
	execDir(cwd, "git", "fetch", "--all")
	formatter, err := createDefaultFormatter(token, config.Repository, config.TicketURL)
	if err != nil {
		return cli.NewExitError(err, 3)
	}
	changelog, nextVersion, err := generateReleaseAndChangelog(cwd, branch, formatter, config)
	if err != nil {
		return cli.NewExitError(err, 4)
	}
	currentChangelog, err := ioutil.ReadFile(changelogfile)
	_, ok := err.(*os.PathError)
	if err != nil && !ok {
		return cli.NewExitError(err, 5)
	}
	if nextVersion == nil {
		return cli.NewExitError(errors.New("could not calculate next version"), 6)
	}
	err = ioutil.WriteFile(changelogfile, []byte(fmt.Sprintf("%s\n\n\n%s", changelog, currentChangelog)), os.ModePerm)
	if err != nil {
		return cli.NewExitError(err, 7)
	}
	err = ioutil.WriteFile(versionPath, []byte(nextVersion.String()), os.ModePerm)
	if err != nil {
		return cli.NewExitError(err, 8)
	}
	return nil
}

func generateReleaseAndChangelog(cwd, branch string, formatter changelog.FormatFunc, config *config.Config) (string, *semver.Version, error) {
	log.Println("generate release..")
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
	commits, err := repo.GetHistoryUntil(latestReleaseCommit.Hash)
	if err != nil {
		return "", nil, err
	}
	if len(commits) == 0 {
		return "", nil, ErrNoCommits
	}
	log.Printf("found %d commits since last release commit", len(commits))
	nextVersion, err := calcReleaseVersion(commits[0], branch, version, config.BranchSuffix, commits.MaxChange())
	if err != nil {
		return "", nil, err
	}
	log.Printf("next version: %s", nextVersion)

	cl := changelog.New(config.Types, formatter)
	changelog := cl.Create(commits, nextVersion)
	return changelog, nextVersion, nil
}
