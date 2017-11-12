package changelog

import (
	"fmt"
	"testing"

	"github.com/Masterminds/semver"

	"github.com/moolen/asdf/repository"
)

var mylog = `## 2.1.3-rc123 (2017-11-12)

#### 

test2
test3

#### Foo

test1
test4

`

func TestChangelog(t *testing.T) {
	typeMap := map[string]string{
		"foo": "Foo",
	}
	cl := New(typeMap, func(commit *repository.Commit) string {
		return fmt.Sprintf("%s\n", commit.Message)
	})
	commits := []*repository.Commit{
		&repository.Commit{
			Message: "test1",
			Type:    "foo",
		},
		&repository.Commit{
			Message: "test2",
			Type:    "",
		},
		&repository.Commit{
			Message: "test3",
			Type:    "",
		},
		&repository.Commit{
			Message: "test4",
			Type:    "foo",
		},
	}
	version := semver.MustParse("2.1.3-rc123")
	log := cl.Create(commits, version)
	if log != mylog {
		t.Fatalf("changelog did not match\nexpected\n%#v\ngot\n%#v\n", mylog, log)
	}
}

func TestTrim(t *testing.T) {
	table := []struct {
		in  string
		out string
	}{
		{
			in:  "1234567890",
			out: "12345678",
		},
		{
			in:  "1234",
			out: "1234",
		},
	}

	for i, r := range table {
		out := TrimSHA(r.in)
		if out != r.out {
			t.Fatalf("[%d] expected %s, got %s", i, r.out, out)
		}
	}
}

func TestDefaultFormatFunc(t *testing.T) {
	table := []struct {
		in  *repository.Commit
		out string
	}{
		{
			in: &repository.Commit{
				Message: "message",
				Hash:    "1234",
			},
			out: "* message (1234) \n",
		},
		{
			in: &repository.Commit{
				Message: "message",
				Hash:    "1234",
				Ticket:  "TEST-123",
			},
			out: "* message [TEST-123] (1234) \n",
		},
	}

	for i, r := range table {
		out := DefaultFormatFunc(r.in)
		if out != r.out {
			t.Fatalf("[%d] expected %#v, got %#v", i, r.out, out)
		}
	}
}

func TestURLFormatFunc(t *testing.T) {

	table := []struct {
		in        *repository.Commit
		formatter FormatFunc
		out       string
	}{
		{
			in: &repository.Commit{
				Message: "message",
				Hash:    "1234",
			},
			formatter: URLFormatFunc("http://example.com/{TICKET_ID}"),
			out:       "* message (1234) \n",
		},
		{
			in: &repository.Commit{
				Message: "message",
				Hash:    "1234",
				Ticket:  "TEST-123",
			},
			formatter: URLFormatFunc("http://example.com/{TICKET_ID}"),
			out:       "* message [TEST-123](http://example.com/TEST-123) (1234) \n",
		},
		{
			in: &repository.Commit{
				Message: "message",
				Hash:    "1234",
				Ticket:  "TEST-123",
			},
			formatter: URLFormatFunc("http://example.com/"),
			out:       "* message [TEST-123](http://example.com/) (1234) \n",
		},
	}

	for i, r := range table {
		out := r.formatter(r.in)
		if out != r.out {
			t.Fatalf("[%d] expected %#v, got %#v", i, r.out, out)
		}
	}
}
