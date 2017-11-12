package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/moolen/asdf/config"
	"github.com/moolen/asdf/fetcher"
	"github.com/moolen/asdf/repository"
)

type FakeFetcher struct{}

func (f FakeFetcher) Fetch() ([]*fetcher.PullRequest, error) {
	return []*fetcher.PullRequest{
		&fetcher.PullRequest{
			ID:     1,
			Title:  "REFPR-1",
			Merged: true,
		},
		&fetcher.PullRequest{
			ID:     2,
			Title:  "REFPR-2",
			Merged: true,
		},
		&fetcher.PullRequest{
			ID:     3,
			Title:  "REFPR-2",
			Merged: true,
		},
		&fetcher.PullRequest{
			ID:     4,
			Title:  "unrelated hotfix",
			Merged: true,
		},
		&fetcher.PullRequest{
			ID:     5,
			Title:  "",
			Merged: true,
		},
	}, nil
}

func TestPrepareRepo(t *testing.T) {

	conf, err := config.FromJSON(strings.NewReader("{}"))
	if err != nil {
		panic(err)
	}

	table := []struct {
		commits []string
		version *semver.Version
		err     error
	}{
		{
			commits: []string{},
			version: nil,
			err:     ErrNoCommits,
		},
		{
			commits: []string{
				"fix(TEST-123): fixing some things",
			},
			version: semver.MustParse("1.0.1"),
		},
		{
			commits: []string{
				"feat(TEST-1): feature 1",
				"feat(TEST-2): feature 2",
				"fix(TEST-123): fixing some things",
			},
			version: semver.MustParse("1.1.0"),
		},
		{
			commits: []string{
				"feat(TEST-2): feature 2",
				"fix(TEST-123): fixing some things",
				"breaking(YALA-123): DEDALDSALD",
				"breaking(YALA-345): 235234",
			},
			version: semver.MustParse("2.0.0"),
		},
	}

	for i, row := range table {
		repo := createRepository()
		for _, commit := range row.commits {
			createAndCommit(repo, commit)
		}
		changelog, nextVersion, err := generateReleaseAndChangelog(repo, "master", &FakeFetcher{}, conf)
		if err != row.err {
			t.Fatalf("[%d]\nexpected %#v\n got %#v", i, row.err, err)
		}
		if !reflect.DeepEqual(nextVersion, row.version) {
			fmt.Println(changelog)
			t.Fatalf("[%d]\nexpected %s\n got %s", i, row.version, nextVersion)
		}
	}
}

func TestPRFormatter(t *testing.T) {
	formatter, err := createPRFormatter(&FakeFetcher{}, "http://example.com/{TICKET_ID}/fart")
	if err != nil {
		t.Fatal(err)
	}
	table := []struct {
		in  *repository.Commit
		out string
	}{
		{
			in:  &repository.Commit{},
			out: "*  () \n",
		},
		{
			in: &repository.Commit{
				Message: "mymessage",
				Ticket:  "UNREFERENCED-1",
			},
			out: "* mymessage [UNREFERENCED-1](http://example.com/UNREFERENCED-1/fart) \n",
		},
		{
			in: &repository.Commit{
				Message: "mymessage",
				Ticket:  "REFPR-1",
			},
			out: "* mymessage [REFPR-1](http://example.com/REFPR-1/fart) (#1) \n",
		},
		{
			in: &repository.Commit{
				Message: "mymessage",
				Ticket:  "REFPR-2",
			},
			out: "* mymessage [REFPR-2](http://example.com/REFPR-2/fart) (#2, #3) \n",
		},
	}
	for i, row := range table {
		out := formatter(row.in)
		if out != row.out {
			t.Fatalf("[%d] expected\n%#v\ngot\n%#v", i, row.out, out)
		}
	}
}
func TestCalcNextVersion(t *testing.T) {
	table := []struct {
		commit      *repository.Commit
		branch      string
		version     *semver.Version
		nextVersion *semver.Version
		suffixMap   map[string]string
		change      config.Change
		err         error
	}{
		{
			commit:      &repository.Commit{},
			branch:      "master",
			version:     semver.MustParse("1.2.3"),
			nextVersion: semver.MustParse("1.2.4"),
			change:      config.ChangePatch,
		},
		{
			commit:      &repository.Commit{},
			branch:      "master",
			version:     semver.MustParse("1.2.3-rc400"),
			nextVersion: semver.MustParse("1.2.3"),
			change:      config.ChangePatch,
		},
		{
			commit: &repository.Commit{
				Hash: "1234",
			},
			branch:      "devrelease",
			version:     semver.MustParse("1.2.3"),
			nextVersion: semver.MustParse("1.2.3-dev1234"),
			suffixMap: map[string]string{
				"devrelease": "dev{COMMIT_SHA}",
			},
			change: config.ChangeMinor,
		},
		{
			commit:      &repository.Commit{},
			branch:      "release",
			version:     semver.MustParse("1.2.3-rc1"),
			nextVersion: semver.MustParse("1.2.3-rc2"),
			suffixMap: map[string]string{
				"release": "rc{RELEASE_NUMBER}",
			},
			change: config.ChangePatch,
		},
		{
			commit:      &repository.Commit{},
			branch:      "beta",
			version:     semver.MustParse("2.0.0-beta.1"),
			nextVersion: semver.MustParse("2.0.0-beta.2"),
			suffixMap: map[string]string{
				"beta": "beta.{RELEASE_NUMBER}",
			},
			change: config.ChangePatch,
		},
	}

	for i, row := range table {
		next, err := calcNextVersion(row.commit, row.branch, row.version, row.suffixMap, row.change)
		if err != row.err {
			t.Fatalf("[%d] expected %s\ngot %s", i, row.err, err)
		}
		if !reflect.DeepEqual(row.nextVersion, next) {
			t.Fatalf("[%d] expected %s\ngot %s", i, row.nextVersion, next)
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
	createAndCommit(repoPath, "initial commit")
	execDir(repoPath, "git", "tag", "1.0.0")
	execDir(repoPath, "git", "push", "origin", "master", "--tags")
	return repoPath
}

func createAndCommit(repo, message string) {
	file, _ := ioutil.TempFile(repo, "")
	file.Close()
	execDir(repo, "git", "add", "-A")
	execDir(repo, "git", "commit", "-m", message)
}

func createVersionFile(repo, version string) {
	err := ioutil.WriteFile(path.Join(repo, "VERSION"), []byte(version), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
