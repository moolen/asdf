package main

import (
	"testing"

	"github.com/moolen/asdf/fetcher"
	"github.com/moolen/asdf/repository"
)

func TestPRFormatter(t *testing.T) {
	formatter, err := createPRFormatter(&FakeFetcher{}, "http://example.com/{SCOPE}/fart")
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
				Subject: "mymessage",
				Scope:   "UNREFERENCED-1",
			},
			out: "* mymessage [UNREFERENCED-1](http://example.com/UNREFERENCED-1/fart) \n",
		},
		{
			in: &repository.Commit{
				Subject: "mymessage",
				Scope:   "REFPR-1",
			},
			out: "* mymessage [REFPR-1](http://example.com/REFPR-1/fart) (#1) \n",
		},
		{
			in: &repository.Commit{
				Subject: "mymessage",
				Scope:   "REFPR-2",
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
