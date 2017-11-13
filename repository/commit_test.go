package repository

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCommitParser(t *testing.T) {
	table := []struct {
		in      string
		err     error
		commits []*Commit
	}{
		{
			in: "2d7ea9249b0afb39c22da7774669738a7e56ff22~Ü>8~#Ä~8<Ü~1591e972ca68d72430ab159100f87683c2080508~Ü>8~#Ä~8<Ü~1510488640~Ü>8~#Ä~8<Ü~Moritz Johner~Ü>8~#Ä~8<Ü~beller.moritz@googlemail.com~Ü>8~#Ä~8<Ü~feat(TEST-2): feature 2",
			commits: []*Commit{
				&Commit{
					ParentHashes: "2d7ea9249b0afb39c22da7774669738a7e56ff22",
					Hash:         "1591e972ca68d72430ab159100f87683c2080508",
					Author: CommitAuthor{
						Name:  "Moritz Johner",
						Email: "beller.moritz@googlemail.com",
					},
					Date:    time.Unix(1510488640, 0),
					Type:    "feat",
					Scope:   "TEST-2",
					Subject: "feature 2",
				},
			},
		},
		{
			in: "2d7ea9249b0afb39c22da7774669738a7e56ff22~Ü>8~#Ä~8<Ü~1591e972ca68d72430ab159100f87683c2080508~Ü>8~#Ä~8<Ü~1510488640~Ü>8~#Ä~8<Ü~Moritz Johner~Ü>8~#Ä~8<Ü~beller.moritz@googlemail.com~Ü>8~#Ä~8<Ü~feat(TEST-2): feature 2\n((((((((----))))))))\nFOOBAR\nBAZLER\n((((((((----))))))))",
			commits: []*Commit{
				&Commit{
					ParentHashes: "2d7ea9249b0afb39c22da7774669738a7e56ff22",
					Hash:         "1591e972ca68d72430ab159100f87683c2080508",
					Author: CommitAuthor{
						Name:  "Moritz Johner",
						Email: "beller.moritz@googlemail.com",
					},
					Date:    time.Unix(1510488640, 0),
					Type:    "feat",
					Scope:   "TEST-2",
					Subject: "feature 2",
					Body:    "FOOBAR\nBAZLER\n",
				},
			},
		},
		{
			// invalid time: fafafafafaffafafasdasd
			in:      "2d7ea9249b0afb39c22da7774669738a7e56ff22~Ü>8~#Ä~8<Ü~1591e972ca68d72430ab159100f87683c2080508~Ü>8~#Ä~8<Ü~fafafafafaffafafasdasd~Ü>8~#Ä~8<Ü~Moritz Johner~Ü>8~#Ä~8<Ü~beller.moritz@googlemail.com~Ü>8~#Ä~8<Ü~feat(TEST-2): feature 2",
			commits: nil,
			err:     strconv.ErrSyntax,
		},
	}

	for i, row := range table {
		commits, err := ParseCommits(strings.NewReader(row.in), DefaultMapFunc)
		if !reflect.DeepEqual(row.err, err) {
			numErr, ok := err.(*strconv.NumError)
			if !ok {
				t.Fatalf("cannot convert to NumErr: %s", err)
			}
			if numErr.Err != row.err {
				t.Fatalf("[%d] expected %#v\ngot %#v\n", i, row.err, err)
			}
		}
		if !reflect.DeepEqual(row.commits, commits) {
			t.Fatalf("[%d] expected \n%#v\ngot \n%#v", i, row.commits[0], commits[0])
		}
	}
}

func TestDefaultMapFunc(t *testing.T) {
	table := []struct {
		in      string
		tp      string
		scope   string
		subject string
	}{
		{
			in:      "feat(TICKK-123): foobar booman",
			tp:      "feat",
			scope:   "TICKK-123",
			subject: "foobar booman",
		},
		{
			in:      "fart(): foobar booman",
			tp:      "fart",
			subject: "foobar booman",
		},
		{
			in:      "(TICKK-123): foobar booman",
			scope:   "TICKK-123",
			subject: "foobar booman",
		},
		{
			in:      "fang: foobar booman",
			tp:      "fang",
			subject: "foobar booman",
		},
		{
			in:      "fang foobar booman",
			subject: "fang foobar booman",
		},
		{
			in:      "feat(1): silly fix we have a maximum line length here. everything >50 chars should be redacted",
			tp:      "feat",
			scope:   "1",
			subject: "silly fix we have a maximum line length here. ever",
		},
	}
	for i, row := range table {
		tp, scope, msg := DefaultMapFunc(row.in)
		if tp != row.tp {
			t.Fatalf("[%d] wrong type: expected %#v, got %#v", i, row.tp, tp)
		}
		if scope != row.scope {
			t.Fatalf("[%d] wrong type: expected %#v, got %#v", i, row.scope, scope)
		}
		if msg != row.subject {
			t.Fatalf("[%d] wrong type: expected %#v, got %#v", i, row.subject, msg)
		}
	}
}
