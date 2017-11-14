package main

import (
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/urfave/cli"
)

const (
	configFilename = "asdf.json"
	flagRevision   = "revision"
	flagDir        = "dir"
	flagFile       = "file"
	flagChangelog  = "changelog"
	flagLatest     = "latest"
	flagVersion    = "version"
)

var errNoRevision = errors.New("revision is required")
var errNoFile = errors.New("file is required")
var errNoVersionProvided = errors.New("no version provided")

// DefaultTypeMap contains a mapping of types to groups
// which are used to render the changelog
var DefaultTypeMap = map[string]string{
	"feat":     "Feature",
	"breaking": "Breaking Changes",
	"fix":      "Bug Fixes",
	"perf":     "Performance Improvements",
	"revert":   "Reverted",
	"docs":     "Documentation",
	"refactor": "Code Refactoring",
	"test":     "Tests",
	"chore":    "Chores",
}

func main() {
	app := cli.NewApp()
	app.Name = "asdf"
	app.Version = "0.1.0"
	app.Usage = "Changelog generation based on semantic commit messages.\n   "
	app.Usage += "The commit messages subject should follow this very convention:\n   "
	app.Usage += "<type>(scope): <description>\n\n   "
	app.Usage += "Example commit messages:\n\n   "
	app.Usage += "feat(TICKET-123): implementing a feature\n   "
	app.Usage += "fix: fixed something\n   "
	app.Usage += "(TICKET-123): some message\n\n   "
	app.Usage += "Only the Commit Subject (first line, 50 characters)\n   "

	app.Commands = []cli.Command{
		{
			Name:    "next",
			Aliases: []string{"n"},
			Usage:   "tells you the next version based on a revision or last file modification",
			Flags:   nextFlags(),
			Action:  nextCommand,
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generates a changelog and the next version based on semantic commits and writes them to file",
			Flags:   generateFlags(),
			Action:  generateCommand,
		},
		{
			Name:    "changelog",
			Aliases: []string{"c"},
			Usage:   "generates only the changelog and writes it to stdout",
			Flags:   changelogFlags(),
			Action:  changelogCommand,
		},
	}
	app.Flags = globalFlags()

	app.Run(os.Args)
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  flagDir,
			Value: "",
			Usage: "set the current wokring directory",
		},
	}
}

func getCwd(c *cli.Context) (string, error) {
	cwd := c.String("dir")
	if cwd == "" {
		dir, err := os.Executable()
		if err != nil {
			return "", err
		}
		cwd = path.Dir(dir)
	}
	return cwd, nil
}

func execDir(dir, cmd string, things ...string) {
	c := exec.Command(cmd, things...)
	c.Dir = dir
	err := c.Run()
	if err != nil {
		panic(err)
	}
}
