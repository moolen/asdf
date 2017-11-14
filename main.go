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
)

// ErrNoRevision the use has to specify a revision(range)
// set man 7 gitrevisions
var ErrNoRevision = errors.New("revision is required")
var ErrNoFile = errors.New("file is required")

func main() {
	app := cli.NewApp()
	app.Name = "asdf"
	app.Version = "0.1.0"
	app.Usage = "Changelog generation based on semantic commit messages.\n   "
	app.Usage += "The changelog generator will ask Github for pull requests that\n   "
	app.Usage += "contain the Ticket ID and will include them in the changelog\n\n   "
	app.Usage += "The commit message subject should follow this very convention:\n   "
	app.Usage += "<type>(scope): <description>\n\n   "
	app.Usage += "Example commit messages:\n\n   "
	app.Usage += "feat(TICKET-123): implementing a feature\n   "
	app.Usage += "fix: fixed something\n   "
	app.Usage += "(TICKET-123): some message\n\n   "
	app.Usage += "Only the Commit Subject (first line, 50 characters)\n   "
	app.Usage += "will be parsed. The tickets will be linked if a URL is set in the configuration file\n   "

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
