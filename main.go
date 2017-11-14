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
	app.Usage = "Changelog and version generation based on semantic commit messages.\n\n   "
	app.Usage += "Specification about the structure is here: http://conventionalcommits.org\n   "
	app.Usage += "All commit messages should follow this very convention:\n   "
	app.Usage += "-------------------------------\n   "
	app.Usage += "<type>(scope): <subject>\n\n   "
	app.Usage += "<body>\n   "
	app.Usage += "-------------------------------\n   "
	app.Usage += "\n   "
	app.Usage += "Example commit message subjects:\n\n   "
	app.Usage += "feat(TICKET-123): implementing a feature\n   "
	app.Usage += "fix: fixed something\n   "
	app.Usage += "(TICKET-123): some message\n\n   "

	app.Commands = []cli.Command{
		{
			Name:    "next-version",
			Aliases: []string{"n"},
			Usage:   "Tells you the next version you want to release. By default it uses a VERSION file to fetch the history since the last release. the file location may be overridden via --" + flagFile,
			Flags:   nextFlags(),
			Action:  nextCommand,
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generates a changelog and the next version based on semantic commits and writes them to files",
			Flags:   generateFlags(),
			Action:  generateCommand,
		},
		{
			Name:    "changelog",
			Aliases: []string{"c"},
			Usage:   "generates the changelog and writes it to stdout. By default it uses a VERSION file to fetch the history since the last release. This can be overridden by defining a --" + flagVersion + " and --" + flagRevision,
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
