package main

import (
	"flag"
	"reflect"
	"testing"

	"github.com/moolen/asdf/repository"
	"github.com/urfave/cli"
)

func TestChangelogCommand(t *testing.T) {
	table := []struct {
		args []string
		err  error
	}{
		{
			args: []string{"--dir"},
			err:  cli.NewExitError(errNoCommits, 5),
		},
		{
			args: []string{"--version", "xyz", "--revision", "HEAD", "--dir"},
			err:  cli.NewExitError(errNoSemverVersion, 2),
		},
		{
			args: []string{"--revision", "", "--version", "1.2.3", "--dir"},
			err:  cli.NewExitError(errNoCommits, 5),
		},
		{
			args: []string{"--revision", "foobar", "--version", "1.2.3", "--dir"},
			err:  cli.NewExitError(repository.ErrExec, 3),
		},
		{
			args: []string{"--revision", "HEAD", "--version", "2.3.4", "--dir"},
			err:  nil,
		},
	}

	for i, row := range table {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		flags := append(changelogFlags(), globalFlags()...)
		for _, flag := range flags {
			flag.Apply(flagSet)
		}
		repo := createRepository()
		flagSet.Parse(append(row.args, repo))
		ctx := cli.NewContext(&cli.App{}, flagSet, nil)
		err := changelogCommand(ctx)
		if !reflect.DeepEqual(err, row.err) {
			t.Fatalf("[%d] expected\n%#v\ngot\n%#v", i, row.err, err)
		}

	}

}
