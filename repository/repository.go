package repository

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
)

// ErrExec is returned if a git command fails.
// That might happen if the command is executed in the wrong directory
// or the git command is not found
var ErrExec = errors.New("git command failed")

// ErrNoHistory is returned when no history was found
// for the specified commit
var ErrNoHistory = errors.New("no history found")

// GitRepository is a interface to a git repository.
// You can access the history, commits and files through this
// struct
type GitRepository struct {
	Path          string
	CommitMapFunc CommitMapFunc
}

// New creates a new Repository
func New(repoPath string, mapFunc CommitMapFunc) *GitRepository {
	return &GitRepository{
		Path:          repoPath,
		CommitMapFunc: mapFunc,
	}
}

// LatestChangeOfFile gives us the commit of the latest change of that file
func (r *GitRepository) LatestChangeOfFile(filename string) (*Commit, error) {
	out, _, err := execDir(r.Path, "git", "log", "-n1", "--format="+logFormatter, "--", filename)
	if err != nil {
		return nil, err
	}
	commits, err := ParseCommits(out, r.CommitMapFunc)
	if err != nil {
		return nil, err
	}
	if len(commits) == 0 {
		return nil, ErrNoHistory
	}
	commit := commits[0]
	return commit, nil
}

// GetHistoryUntil returns all commits from HEAD to the specified commit
func (r *GitRepository) GetHistoryUntil(revision string) (Commits, error) {
	var commits Commits
	out, _, err := execDir(r.Path, "git", "log", "--format="+logFormatter, revision+"..HEAD")
	if err != nil {
		return commits, err
	}
	return ParseCommits(out, r.CommitMapFunc)
}

// GetHistory returns all commits defined by a gitrevision
// Examples:
// - "13c2a8c..HEAD"
// - "develop..master"
// - "HEAD^1"
// For further information read `man 7 gitrevisions`
func (r *GitRepository) GetHistory(gitrevisions string) (Commits, error) {
	var commits Commits
	out, _, err := execDir(r.Path, "git", "log", "--format="+logFormatter, gitrevisions)
	if err != nil {
		return commits, err
	}
	return ParseCommits(out, r.CommitMapFunc)
}

// execDir executes a command in a specific directory
func execDir(dir, cmd string, things ...string) (io.Reader, io.Reader, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c := exec.Command(cmd, things...)
	c.Dir = dir
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		return nil, nil, ErrExec
	}
	return &stdout, &stderr, nil
}
