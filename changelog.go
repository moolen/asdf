package main

import (
	"os"
	"path"

	"github.com/Masterminds/semver"
	"github.com/moolen/asdf/changelog"
	"github.com/moolen/asdf/repository"
	"github.com/urfave/cli"
)

// changelog is a stateless command that, given a range,
// will write the changelog to stdout
func changelogCommand(c *cli.Context) error {
	var err error
	var commits repository.Commits
	var version *semver.Version
	revision := c.String(flagRevision)
	versionString := c.String(flagVersion)
	versionFile := c.String(flagFile)
	cwd, err := getCwd(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	versionPath := path.Join(cwd, versionFile)
	repo := repository.New(cwd, repository.DefaultMapFunc)

	// 2nd use-case: supply revision + version explicitly
	if revision != "" && versionString != "" {
		version, err = semver.NewVersion(versionString)
		if err != nil {
			return cli.NewExitError(errNoSemverVersion, 2)
		}
		commits, err = repo.GetHistory(revision)
		if err != nil {
			return cli.NewExitError(err, 3)
		}
	} else if versionFile != "" {
		version, err = readVersionFile(versionPath)
		if err != nil && !os.IsNotExist(err) {
			return cli.NewExitError(err, 4)
		}
		commit, err := repo.LatestChangeOfFile(versionPath)
		if err != nil {
			return cli.NewExitError(err, 5)
		}
		commits, err = repo.GetHistoryUntil(commit.Hash)
		if err != nil {
			return cli.NewExitError(err, 6)
		}
	}

	if len(commits) == 0 {
		return cli.NewExitError(errNoCommits, 5)
	}
	cl := changelog.New(DefaultTypeMap, changelog.DefaultFormatFunc)
	os.Stdout.WriteString(cl.Create(commits, version))
	return nil
}

func changelogFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  flagRevision,
			Usage: "revision to calculate the diff. Works only together with --" + flagVersion,
		},
		cli.StringFlag{
			Name:  flagVersion,
			Value: "",
			Usage: "set the release version explicitly works only in conjunction with --" + flagRevision,
		},
		cli.StringFlag{
			Name:  flagFile,
			Value: "VERSION",
			Usage: "file to use to get the commit of last modification. That file must include the latest version",
		},
	}
}
