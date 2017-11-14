package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/moolen/asdf/changelog"
	"github.com/urfave/cli"
)

func TestPrepareRepo(t *testing.T) {

	table := []struct {
		commits map[string]string
		version *semver.Version
		err     error
	}{
		{
			commits: map[string]string{},
			version: nil,
			err:     errNoCommits,
		},
		{
			commits: map[string]string{
				"fix(TEST-123): fixing some things": "",
			},
			version: semver.MustParse("1.0.1"),
		},
		{
			commits: map[string]string{
				"feat(TEST-1): feature 1":           "",
				"feat(TEST-2): feature 2":           "",
				"fix(TEST-123): fixing some things": "",
			},
			version: semver.MustParse("1.1.0"),
		},
		{
			commits: map[string]string{
				"feat(TEST-2): feature 2":           "",
				"fix(TEST-123): fixing some things": "",
				"breaking(YALA-123): DEDALDSALD":    "",
				"breaking(YALA-345): 235234":        "",
			},
			version: semver.MustParse("1.1.0"),
		},
		{
			commits: map[string]string{
				"fix(NOTREALLY): yolo": "BREAKING CHANGE: your mom",
			},
			version: semver.MustParse("2.0.0"),
		},
	}

	for i, row := range table {
		repo := createRepository()
		for subject, body := range row.commits {
			createAndCommit(repo, subject, body)
		}
		fmt.Printf("%#v", path.Join(repo, "VERSION"))
		changelog, nextVersion, err := generateReleaseAndChangelog(repo, "VERSION", changelog.DefaultFormatFunc)
		if err != row.err {
			t.Fatalf("[%d]\nexpected %#v\n got %#v", i, row.err, err)
		}
		if !reflect.DeepEqual(nextVersion, row.version) {
			fmt.Println(changelog)
			t.Fatalf("[%d]\nexpected %s\n got %s", i, row.version, nextVersion)
		}
	}
}

func TestGenerateCommand(t *testing.T) {
	table := []struct {
		commits map[string]string
		args    []string
		err     error
	}{
		{
			args: []string{"--dir"},
			err:  cli.NewExitError(errNoCommits, 4),
		},
		{
			commits: map[string]string{
				"feat: foobar": "",
			},
			args: []string{"--dir"},
			err:  nil,
		},
	}

	for i, row := range table {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		flags := append(generateFlags(), globalFlags()...)
		for _, flag := range flags {
			flag.Apply(flagSet)
		}
		repo := createRepository()
		for subject, body := range row.commits {
			createAndCommit(repo, subject, body)
		}
		flagSet.Parse(append(row.args, repo))
		ctx := cli.NewContext(&cli.App{}, flagSet, nil)
		err := generateCommand(ctx)
		if !reflect.DeepEqual(err, row.err) {
			t.Fatalf("[%d] expected\n%#v\ngot\n%#v", i, row.err, err)
		}

	}
}

// createRepository gives us a git repository
// with one single commit that contains a VERSION file and a tag `1.0.0`.
// Those changes are reflected at the remote bare repository
func createRepository() string {
	repoPath, _ := ioutil.TempDir("", "asdf")
	bareRepoPath, _ := ioutil.TempDir("", "asdf")
	execDir(repoPath, "git", "init")
	execDir(bareRepoPath, "git", "init", "--bare")
	execDir(repoPath, "git", "remote", "add", "origin", bareRepoPath)

	createVersionFile(repoPath, "1.0.0")
	createAndCommit(repoPath, "initial commit", "")
	execDir(repoPath, "git", "tag", "1.0.0")
	execDir(repoPath, "git", "push", "origin", "master", "--tags")
	return repoPath
}

func createAndCommit(repo, subject, body string) {
	file, _ := ioutil.TempFile(repo, "")
	file.Close()
	execDir(repo, "git", "add", "-A")
	execDir(repo, "git", "commit", "-m", subject, "-m", body)
}

func createVersionFile(repo, version string) {
	err := ioutil.WriteFile(path.Join(repo, "VERSION"), []byte(version), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
