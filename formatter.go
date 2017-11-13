package main

import (
	"fmt"
	"strings"

	"github.com/moolen/asdf/changelog"
	"github.com/moolen/asdf/fetcher"
	"github.com/moolen/asdf/repository"
)

func createDefaultFormatter(token, githubRepo, ticketURL string) (changelog.FormatFunc, error) {
	formatter := changelog.DefaultFormatFunc
	// use github PRs
	if token != "" && githubRepo != "" {
		fetch, err := fetcher.New(token, githubRepo)
		if err != nil {
			return formatter, err
		}
		formatter, err = createPRFormatter(fetch, ticketURL)
		if err != nil {
			return formatter, err
		}
	}
	return formatter, nil
}

// this returns a FormatFunc for the changelog.
// it will check if PRs matches the Ticket ID from the commit
// and renders the PR ID in the changelog
func createPRFormatter(fetcher fetcher.PullRequestFetcher, url string) (changelog.FormatFunc, error) {
	pullRequests, err := fetcher.Fetch()
	PullRequestMap := make(map[string][]string)
	if err != nil {
		return nil, err
	}
	formatPullRequestID := func(ID int) string {
		return fmt.Sprintf("#%d", ID)
	}
	// parse all pull requests and put them into a map
	// so we can have a easy direct lookup
	for _, pr := range pullRequests {
		if pr.Merged == true && pr.Title != "" {
			matched := pullRequestTitleRegex.FindAllStringSubmatch(pr.Title, 1)
			if len(matched) > 0 {
				ticketID := matched[0][0]
				if PullRequestMap[ticketID] == nil {
					PullRequestMap[ticketID] = []string{formatPullRequestID(pr.ID)}
				} else {
					PullRequestMap[ticketID] = append(PullRequestMap[ticketID], formatPullRequestID(pr.ID))
				}
			}
		}
	}
	// return the changelog.FormatFunc
	return func(c *repository.Commit) string {
		if c.Scope != "" {
			var prList string
			ticketURL := strings.Replace(url, "{SCOPE}", c.Scope, -1)
			if len(PullRequestMap[c.Scope]) > 0 {
				prList = strings.Join(PullRequestMap[c.Scope], ", ")
				return fmt.Sprintf("* %s [%s](%s) (%s) \n", c.Subject, c.Scope, ticketURL, prList)
			}
			return fmt.Sprintf("* %s [%s](%s) \n", c.Subject, c.Scope, ticketURL)
		}
		return changelog.DefaultFormatFunc(c)
	}, nil
}
