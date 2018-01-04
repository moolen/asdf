package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/Masterminds/semver"
	"github.com/figome/semantic-changelog/changelog"
	"github.com/figome/semantic-changelog/repository"
	"github.com/urfave/cli"
)

var errNoCommits = errors.New("there is nothing to release: no new commits found")
var errNoSemverVersion = errors.New("version file does not contain a semver version")

// ReleaseToken is repleaced with the prerelease number
// If there was no previous release it will starting with 1
var ReleaseToken = "{RELEASE_NUMBER}"

// CommitToken is replaced within a release and contains the short commit hash
var CommitToken = "{COMMIT_SHA}"

// generateCommand is a stateful operation
// It looks for a VERSION file, calculates the
// changelog based on the commits since this file has changed
func generateCommand(c *cli.Context) error {
	versionFile := c.String(flagFile)
	changelogFile := c.String(flagChangelog)
	cwd, err := getCwd(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	log.Infof("working in dir: %s", cwd)
	versionPath := path.Join(cwd, versionFile)
	changelogfile := path.Join(cwd, changelogFile)
	execDir(cwd, "git", "fetch", "--all")
	changelog, nextVersion, err := generateReleaseAndChangelog(cwd, versionFile, changelog.DefaultFormatFunc)
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

func generateReleaseAndChangelog(cwd, versionfile string, formatter changelog.FormatFunc) (string, *semver.Version, error) {
	versionPath := path.Join(cwd, versionfile)
	versionFile, err := os.Open(versionPath)
	defer versionFile.Close()
	if err != nil {
		return "", nil, errors.New("version file does not exist, please create one")
	}
	version, err := readVersion(versionFile)
	if err != nil {
		return "", nil, errNoSemverVersion
	}
	log.Infof("found version: %s", version)
	repo := repository.New(cwd, repository.DefaultMapFunc)
	latestReleaseCommit, err := repo.LatestChangeOfFile(path.Base(versionPath))
	if err != nil {
		return "", nil, err
	}
	log.Infof("latest release commit: (%s) %s", latestReleaseCommit.Hash, latestReleaseCommit.Subject)
	commits, err := repo.GetHistoryUntil(latestReleaseCommit.Hash)
	if err != nil {
		return "", nil, err
	}
	if len(commits) == 0 {
		return "", nil, errNoCommits
	}
	log.Infof("found %d commits since last release commit", len(commits))
	nextVersion := nextReleaseByChange(version, commits.MaxChange())
	if err != nil {
		return "", nil, err
	}
	log.Infof("next version: %s", nextVersion.String())

	cl := changelog.New(DefaultTypeMap, formatter)
	changelog := cl.Create(commits, &nextVersion)
	return changelog, &nextVersion, nil
}

func generateFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  flagFile,
			Value: "VERSION",
			Usage: "file that holds the version information",
		},
		cli.StringFlag{
			Name:  flagChangelog,
			Value: "CHANGELOG.md",
			Usage: "file that holds the changelog",
		},
	}
}
