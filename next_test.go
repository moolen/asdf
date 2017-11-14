package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/urfave/cli"
)

func TestNextCommand(t *testing.T) {
	table := []struct {
		commits map[string]string
		stdout  string
		args    []string
		err     error
	}{
		{
			args:   []string{"--dir"},
			stdout: "",
			err:    cli.NewExitError(errNoCommits, 6),
		},
		{
			commits: map[string]string{
				"feat: bar": "BREAKING CHANGES: yolo",
			},
			args:   []string{"--dir"},
			stdout: "2.0.0",
			err:    nil,
		},
		{
			commits: map[string]string{
				"feat: bar": "yolo",
			},
			args:   []string{"--dir"},
			stdout: "1.1.0",
			err:    nil,
		},
		{
			commits: map[string]string{
				"fix: bar":  "yolo",
				"feat: bar": "test",
			},
			args:   []string{"--file", "MYVERSIONFILE", "--dir"},
			stdout: "13.15.0",
			err:    nil,
		},
	}

	for i, row := range table {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		flags := append(nextFlags(), globalFlags()...)
		for _, flag := range flags {
			flag.Apply(flagSet)
		}
		repo := createRepository()
		ioutil.WriteFile(path.Join(repo, "MYVERSIONFILE"), []byte("13.14.15"), os.ModePerm)
		for subject, body := range row.commits {
			createAndCommit(repo, subject, body)
		}

		flagSet.Parse(append(row.args, repo))
		ctx := cli.NewContext(&cli.App{}, flagSet, nil)
		stdout := os.Stdout
		tempfile, _ := ioutil.TempFile("", "")
		defer tempfile.Close()
		os.Stdout = tempfile
		err := nextCommand(ctx)
		bytes, blen := ioutil.ReadFile(tempfile.Name())
		os.Stdout = stdout
		if string(bytes) != row.stdout {
			t.Fatalf("[%d] %#v expected\n%#v\ngot\n%#v", i, blen, row.stdout, string(bytes))
		}
		if !reflect.DeepEqual(err, row.err) {
			t.Fatalf("[%d] expected\n%#v\ngot\n%#v", i, row.err, err)
		}

	}
}
