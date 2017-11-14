package main

import (
	"os"

	"github.com/Masterminds/semver"

	"github.com/moolen/asdf/changelog"
	"github.com/moolen/asdf/repository"
	"github.com/urfave/cli"
)

// changelog is a stateless command that, given a range,
// will write the changelog to stdout
func changelogCommand(c *cli.Context) error {
	revision := c.String(flagRevision)
	versionString := c.String(flagVersion)
	if versionString == "" {
		return cli.NewExitError(errNoVersionProvided, 1)
	}
	version, err := semver.NewVersion(versionString)
	if err != nil {
		return cli.NewExitError(errNoSemverVersion, 1)
	}
	if revision == "" {
		return cli.NewExitError(errNoRevision, 2)
	}
	cwd, err := getCwd(c)
	if err != nil {
		return cli.NewExitError(err, 3)
	}
	repo := repository.New(cwd, repository.DefaultMapFunc)
	commits, err := repo.GetHistory(revision)
	if err != nil {
		return cli.NewExitError(err, 4)
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
			Usage: "revision to calculate the diff",
		},
		cli.StringFlag{
			Name:  flagVersion,
			Usage: "the release version",
		},
	}
}
