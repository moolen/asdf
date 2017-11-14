package main

import (
	"os"
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/moolen/asdf/repository"
	"github.com/urfave/cli"
)

func nextCommand(c *cli.Context) error {
	var commit *repository.Commit
	file := c.String(flagFile)
	cwd, err := getCwd(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if file == "" {
		return cli.NewExitError(errNoFile, 2)
	}
	file = path.Join(cwd, file)
	repo := repository.New(cwd, repository.DefaultMapFunc)
	commit, err = repo.LatestChangeOfFile(file)
	if err != nil {
		return cli.NewExitError(err, 3)
	}
	log.Printf("file %s had last change at %s in commit %s", file, commit.Date.Format("2006-01-02"), commit.Hash)
	latest, err := readVersionFile(file)
	log.Printf("found version: %s", latest)
	if err != nil {
		return cli.NewExitError(err, 4)
	}
	commits, err := repo.GetHistoryUntil(commit.Hash)
	if err != nil {
		return cli.NewExitError(err, 5)
	}
	if len(commits) == 0 {
		return cli.NewExitError(errNoCommits, 6)
	}
	log.Printf("commits since last change: %d", len(commits))

	log.Printf("found max change: %s", commits.MaxChange())
	nextVersion := nextReleaseByChange(latest, commits.MaxChange())
	os.Stdout.WriteString(nextVersion.String())
	return nil
}

func nextFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  flagFile,
			Value: "VERSION",
			Usage: "file to use to get the commit of last modification. That file must include the latest version",
		},
	}
}
