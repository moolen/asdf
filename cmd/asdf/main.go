package main

import (
	"errors"
	"os"
	"path"

	"github.com/moolen/asdf/config"
	"github.com/moolen/asdf/fetcher"
	"github.com/urfave/cli"
)

const (
	configFilename = "asdf.json"
)

func main() {
	app := cli.NewApp()
	app.Name = "asdf"
	app.Version = "0.1.0"
	app.Usage = "Changelog generation based on semantic commit messages"
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generates a changelog and the next version based on semantic commits and writes them to file",
			Action: func(c *cli.Context) error {
				cwd, err := getCwd(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				token := c.GlobalString("token")
				if token == "" {
					return cli.NewExitError(errors.New("github token is missing"), 1)
				}
				config, err := config.FromFile(path.Join(cwd, configFilename))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				err = generateRelease(cwd, token, c.GlobalString("branch"), config)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		{
			Name:    "changelog",
			Aliases: []string{"c"},
			Usage:   "generates only the changelog and writes it to stdout",
			Action: func(c *cli.Context) error {
				cwd, err := getCwd(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				token := c.GlobalString("token")
				if token == "" {
					return cli.NewExitError(errors.New("github token is missing"), 1)
				}
				config, err := config.FromFile(path.Join(cwd, configFilename))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				fetcher, err := fetcher.New(token, config.Repository)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				changelog, _, err := generateReleaseAndChangelog(cwd, c.GlobalString("branch"), fetcher, config)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				os.Stdout.WriteString(changelog)
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "branch",
			Value:  "master",
			Usage:  "name of the current branch",
			EnvVar: "RELEASE_BRANCH",
		},
		cli.StringFlag{
			Name:   "token",
			Value:  "",
			Usage:  "github token",
			EnvVar: "RELEASE_GITHUB_TOKEN",
		},
	}

	app.Run(os.Args)
}

func getCwd(c *cli.Context) (string, error) {
	cwd := c.String("dir")
	if cwd == "" {
		dir, err := os.Executable()
		if err != nil {
			return "", cli.NewExitError(err, 1)
		}
		cwd = path.Dir(dir)
	}
	return cwd, nil
}
