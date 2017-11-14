package repository

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestLatestChangeOfFile(t *testing.T) {
	repoPath := createRepository()
	repo := New(repoPath, DefaultMapFunc)
	commit, err := repo.LatestChangeOfFile("VERSION")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if commit.Subject != "initial commit" {
		t.Fatalf("commit message is wrong: expected %s, got %s", "initial commit", commit.Subject)
	}
}

func TestLatestChangeOfFileFail(t *testing.T) {
	repoPath, _ := ioutil.TempDir("", "bh")
	repo := New(repoPath, DefaultMapFunc)
	commit, err := repo.LatestChangeOfFile("VERSION")
	if err != ErrExec {
		t.Fatalf("expected ErrExec, got: %v", err)
	}
	if commit != nil {
		t.Fatalf("commit is wrong: expected nil, got %#v", commit)
	}
}

func TestLatestChangeOfFileNoHistory(t *testing.T) {
	repoPath := createRepository()
	repo := New(repoPath, DefaultMapFunc)
	commit, err := repo.LatestChangeOfFile("DOESNOTEXIST")
	if err != ErrNoHistory {
		t.Fatalf("expected ErrNoHistory, got: %v", err)
	}
	if commit != nil {
		t.Fatalf("commit is wrong: expected nil, got %#v", commit)
	}
}

func TestGetHistory(t *testing.T) {
	repoPath := createRepository()
	repo := New(repoPath, DefaultMapFunc)
	commit, err := repo.LatestChangeOfFile("VERSION")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if commit.Subject != "initial commit" {
		t.Fatalf("commit message is wrong: expected %s, got %s", "initial commit", commit.Subject)
	}

	// add some commits on top
	createAndCommit(repoPath, "first")
	createAndCommit(repoPath, "second")
	createAndCommit(repoPath, "third")
	commits, err := repo.GetHistoryUntil(commit.Hash)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if commits[0].Subject != "third" {
		t.Fatalf("third commit message wrong: expected third, got: %s", commits[0].Subject)
	}
	if commits[1].Subject != "second" {
		t.Fatalf("second commit message wrong: expected second, got: %s", commits[1].Subject)
	}
	if commits[2].Subject != "first" {
		t.Fatalf("first commit message wrong: expected first, got: %s", commits[2].Subject)
	}
}

func TestGetHistoryExecErr(t *testing.T) {
	repoPath, _ := ioutil.TempDir("", "bads")
	repo := New(repoPath, DefaultMapFunc)
	commits, err := repo.GetHistoryUntil("")
	if err != ErrExec {
		t.Fatalf("expected ErrExec, got: %v", err)
	}
	if commits != nil {
		t.Fatalf("unexpected commits found: %#v", commits)
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
