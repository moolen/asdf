package repository

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Commit holds all the relevant information about
// a git commit
type Commit struct {
	Hash         string
	ParentHashes string
	Author       CommitAuthor
	Date         time.Time
	Type         string
	Ticket       string
	Message      string
}

// CommitAuthor holds information regarding the author of the commit
type CommitAuthor struct {
	Name  string
	Email string
}

var commitPattern = regexp.MustCompile("^(\\w*)(?:\\((.*)\\))?\\: (.*)$")

// ErrParse happens, if git log gives us wrong output
var ErrParse = errors.New("could not parse log output")

// CommitMapFunc is a function that receives a commit message,
// parses it, and returns:
// - the type of change (feat/ix/breaking)
// - the ticket ID
// - the stripped message
type CommitMapFunc func(msg string) (commitType string, commitTicket string, commitMessage string)

// used as delimiter to split the values from git log
var delimiter = "~Ü>8~#Ä~8<Ü~"

// see https://git-scm.com/docs/pretty-formats
var formatString = []string{
	"%P",  // parent hashes
	"%H",  // commit hash
	"%at", // author date, UNIX timestamp
	"%an", // author name
	"%ae", // author email
	"%s",  // commit message subject
}

var logFormatter = strings.Join(formatString, delimiter)

// DefaultMapFunc parses the commit message
// and returns a type
func DefaultMapFunc(msg string) (commitType string, commitTicket string, commitMessage string) {
	lines := strings.Split(msg, "\n")
	found := commitPattern.FindAllStringSubmatch(lines[0], -1)
	if len(found) < 1 {
		return "", "", msg
	}
	commitType = strings.ToLower(found[0][1])
	commitTicket = strings.ToUpper(found[0][2])
	commitMessage = fmt.Sprintf("%.50s", strings.ToLower(found[0][3]))
	return
}

// ParseCommits parses a commit message from an io.Reader
// and returns a Commit
func ParseCommits(stdout io.Reader, mapFunc CommitMapFunc) ([]*Commit, error) {
	var commits []*Commit
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, delimiter)
		if len(splitLine) != 6 {
			return nil, ErrParse
		}
		unixSeconds, err := strconv.ParseInt(splitLine[2], 10, 64)
		if err != nil {
			return nil, err
		}
		changedDate := time.Unix(unixSeconds, 0)
		commitType, commitTicket, commitMessage := mapFunc(splitLine[5])
		commits = append(commits, &Commit{
			ParentHashes: splitLine[0],
			Hash:         splitLine[1],
			Date:         changedDate,
			Message:      commitMessage,
			Author: CommitAuthor{
				Name:  splitLine[3],
				Email: splitLine[4],
			},
			Ticket: commitTicket,
			Type:   commitType,
		})
	}
	return commits, nil
}
