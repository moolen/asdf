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
		err  *cli.ExitError
	}{
		{
			args: []string{""},
			err:  cli.NewExitError(ErrNoRevision, 1),
		},
		{
			args: []string{"--revision", "foobar"},
			err:  cli.NewExitError(repository.ErrExec, 3),
		},
		{
			args: []string{"--revision", "HEAD", "--dir"},
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
