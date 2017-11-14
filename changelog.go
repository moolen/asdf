package main

import (
	"os"

	"github.com/moolen/asdf/changelog"
	"github.com/moolen/asdf/config"
	"github.com/moolen/asdf/repository"
	"github.com/urfave/cli"
)

// changelog is a stateless command that, given a range,
// will write the changelog to stdout
func changelogCommand(c *cli.Context) error {
	revision := c.String(flagRevision)
	if revision == "" {
		return cli.NewExitError(ErrNoRevision, 1)
	}
	cwd, err := getCwd(c)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	repo := repository.New(cwd, repository.DefaultMapFunc)
	commits, err := repo.GetHistory(revision)
	if err != nil {
		return cli.NewExitError(err, 3)
	}
	if len(commits) == 0 {
		return cli.NewExitError(ErrNoCommits, 4)
	}
	cl := changelog.New(config.DefaultTypeMap, changelog.DefaultFormatFunc)
	os.Stdout.WriteString(cl.Create(commits, nil))
	return nil
}

func changelogFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  flagRevision,
			Usage: "revision to calculate the diff",
		},
	}
}
