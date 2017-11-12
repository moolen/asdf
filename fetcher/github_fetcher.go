package fetcher

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// PullRequestFetcher is a interface that provides
// access to pull requests
type PullRequestFetcher interface {
	Fetch() ([]*PullRequest, error)
}

// PullRequest encloses information about a pull request
type PullRequest struct {
	ID     int
	Title  string
	Merged bool
}

// GithubFetcher is a PullRequestFetcher that fetches
// the PRs from the GitHub API
type GithubFetcher struct {
	Owner  string
	Repo   string
	client *github.Client
}

// Fetch returns an array of pull requests
func (g *GithubFetcher) Fetch() ([]*PullRequest, error) {
	var pullRequests []*PullRequest
	PRService := g.client.GetPullRequests()
	PRList, _, err := PRService.List(context.TODO(), "figome", "banking", &github.PullRequestListOptions{
		State: "closed",
	})
	if err != nil {
		return nil, err
	}
	for _, pullRequest := range PRList {
		if pullRequest.ID != nil && pullRequest.Title != nil && pullRequest.Merged != nil {
			log.Printf("pr: %#v", pullRequest)
			pullRequests = append(pullRequests, &PullRequest{
				ID:     *pullRequest.ID,
				Title:  *pullRequest.Title,
				Merged: *pullRequest.Merged,
			})
		}
	}
	return pullRequests, nil
}

// New returns a GithubFetcher and maybe an error
// if the repoSlug is in the wrong format
// repoSlug should be e.g. `moolen/asdf``
func New(token, repoSlug string) (*GithubFetcher, error) {
	split := strings.Split(repoSlug, "/")
	if len(split) != 2 {
		return nil, fmt.Errorf("could not parse: %s", repoSlug)
	}
	return &GithubFetcher{
		Owner: split[0],
		Repo:  split[1],
		client: github.NewClient(oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))),
	}, nil
}
