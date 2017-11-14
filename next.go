package main

import (
	"log"
	"os"

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
		return cli.NewExitError(ErrNoFile, 2)
	}
	repo := repository.New(cwd, repository.DefaultMapFunc)
	commit, err = repo.LatestChangeOfFile(file)
	if err != nil {
		return cli.NewExitError(err, 3)
	}
	log.Printf("file %s had last change at %s in commit %s", file, commit.Date.Format("2006-01-02"), commit.Hash)
	latest, err := readVersionFile(file)
	log.Printf("found version: %s", latest)
	if err != nil {
		return cli.NewExitError(err, 5)
	}
	commits, err := repo.GetHistoryUntil(commit.Hash)
	if err != nil {
		return cli.NewExitError(err, 4)
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
			Usage: "file to use to get the commit of last modification. That file must include the latest version",
		},
	}
}
