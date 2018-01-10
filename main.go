package main

import (
	"errors"
	"os"
	"os/exec"
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
)

const (
	flagRevision  = "revision"
	flagDir       = "dir"
	flagFile      = "file"
	flagPrerelease = "prerelease"
	flagMetadata = "metadata"
	flagChangelog = "changelog"
	flagLatest    = "latest"
	flagVersion   = "version"
	flagDebug     = "debug"
)

var smclVersion string
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
	log.SetFormatter(&log.TextFormatter{})
	app := cli.NewApp()
	app.Name = "smcl"
	app.Version = smclVersion
	app.Usage = "Changelog and version generation based on semantic commit messages.\n\n   "
	app.Usage += "Specification about the structure is here:\n   "
	app.Usage += "https://github.com/figome/figo-rfc/blob/master/docs/COMMIT_MESSAGE.md"

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

	app.Before = func(c *cli.Context) error {
		if c.GlobalBool(flagDebug) {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Flags = globalFlags()

	app.Run(os.Args)
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  flagDir,
			Value: "",
			Usage: "set the current working directory",
		},
		cli.BoolFlag{
			Name:  flagDebug,
			Usage: "show debug logs",
		},
	}
}

func getCwd(c *cli.Context) (string, error) {
	cwd := c.GlobalString("dir")
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
