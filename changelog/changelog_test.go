package changelog

import (
	"fmt"
	"testing"
	"time"

	"github.com/Masterminds/semver"

	"github.com/figome/semantic-changelog/repository"
)

var mylog = `## 2.1.3-rc123 (%s)

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
		return fmt.Sprintf("%s\n", commit.Subject)
	})
	commits := []*repository.Commit{
		{
			Subject: "test1",
			Type:    "foo",
		},
		{
			Subject: "test2",
			Type:    "",
		},
		{
			Subject: "test3",
			Type:    "",
		},
		{
			Subject: "test4",
			Type:    "foo",
		},
	}
	version := semver.MustParse("2.1.3-rc123")
	log := cl.Create(commits, version)
	formattedLog := fmt.Sprintf(mylog, time.Now().Format("2006-01-02"))
	if log != formattedLog {
		t.Fatalf("changelog did not match\nexpected\n%#v\ngot\n%#v\n", formattedLog, log)
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
				Subject: "message",
				Hash:    "1234",
			},
			out: "* message (1234) \n",
		},
		{
			in: &repository.Commit{
				Subject: "message",
				Hash:    "1234",
				Scope:   "TEST-123",
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
