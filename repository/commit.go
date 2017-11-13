package repository

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Change int

const (
	PatchChange Change = iota
	MinorChange
	MajorChange
)

// Commit holds all the relevant information about
// a git commit
type Commit struct {
	Hash         string
	ParentHashes string
	Author       CommitAuthor
	Date         time.Time
	Type         string
	Scope        string
	Subject      string
	Body         string
	Change       Change
}

// CommitAuthor holds information regarding the author of the commit
type CommitAuthor struct {
	Name  string
	Email string
}

// Commits is just a simple list of commits
// that provides convenient functionality
type Commits []*Commit

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

var bodyBeginSeperator = "((((((((----))))))))"
var bodyEndSeperator = "((((((((^^^^))))))))"

var logFormatter = strings.Join(formatString, delimiter) + "%n" + bodyBeginSeperator + "%n%b%n" + bodyEndSeperator

// DefaultMapFunc parses the commit message
// and returns a type
func DefaultMapFunc(msg string) (commitType string, commitScope string, commitMessage string) {
	lines := strings.Split(msg, "\n")
	found := commitPattern.FindAllStringSubmatch(lines[0], -1)
	if len(found) < 1 {
		return "", "", msg
	}
	commitType = strings.ToLower(found[0][1])
	commitScope = strings.ToUpper(found[0][2])
	commitMessage = fmt.Sprintf("%.50s", strings.ToLower(found[0][3]))
	return
}

// ParseCommits parses a commit message from an io.Reader
// and returns a Commit
func ParseCommits(stdout io.Reader, mapFunc CommitMapFunc) ([]*Commit, error) {
	var commits []*Commit
	change := PatchChange
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// first line is always the commit metadata
		metadata := scanner.Text()
		var body string

		// the "rest" is the commit body
		// that contains newline characters
	BodyLoop:
		for scanner.Scan() {
			bodyLine := scanner.Text()
			if bodyLine == bodyBeginSeperator {
				continue BodyLoop
			}
			if bodyLine == bodyEndSeperator {
				break BodyLoop
			}
			body += fmt.Sprintf("%s\n", bodyLine)
		}
		parsedMetadata := strings.Split(metadata, delimiter)
		if len(parsedMetadata) != 6 {
			log.Printf("%#v", parsedMetadata)
			return nil, ErrParse
		}
		unixSeconds, err := strconv.ParseInt(parsedMetadata[2], 10, 64)
		if err != nil {
			return nil, err
		}
		changedDate := time.Unix(unixSeconds, 0)
		commitType, commitScope, commitMessage := mapFunc(parsedMetadata[5])
		if commitType == "feat" {
			change = MinorChange
		}
		if strings.HasPrefix(body, "BREAKING CHANGE") {
			change = MajorChange
		}
		commits = append(commits, &Commit{
			ParentHashes: parsedMetadata[0],
			Hash:         parsedMetadata[1],
			Date:         changedDate,
			Subject:      commitMessage,
			Body:         body,
			Author: CommitAuthor{
				Name:  parsedMetadata[3],
				Email: parsedMetadata[4],
			},
			Change: change,
			Scope:  commitScope,
			Type:   commitType,
		})
	}
	return commits, nil
}

// MaxChange gives us the max
func (commits Commits) MaxChange() Change {
	max := PatchChange
	for _, commit := range commits {
		if max < commit.Change {
			max = commit.Change
		}
	}
	return max
}
