package changelog

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/moolen/asdf/repository"
)

// FormatFunc is called for every commit and should
// return a pretty formatted commit message
// which will be used for the changelog generation
type FormatFunc func(*repository.Commit) string

// Changelog is used to create pretty changelog documents
// provide a TypeMap to group commit messages
// or a FormatFunc to style the messages
type Changelog struct {
	TypeMap    map[string]string
	FormatFunc FormatFunc
}

// New creates a new Changelog struct
func New(typeMap map[string]string, format FormatFunc) *Changelog {
	return &Changelog{
		TypeMap:    typeMap,
		FormatFunc: format,
	}
}

// Create returns a pretty changelog as a string given an array of commits
// This uses the TypeMap to group the commits by type and
// formats every commit with the FormatFunc
func (c *Changelog) Create(commits []*repository.Commit, newVersion *semver.Version) string {
	ret := fmt.Sprintf("## %s (%s)\n\n", newVersion.String(), time.Now().UTC().Format("2006-01-02"))
	typeScopeMap := make(map[string]string)
	for _, commit := range commits {
		typeScopeMap[commit.Type] += c.FormatFunc(commit)
	}
	for _, t := range getSortedKeys(&typeScopeMap) {
		msg := typeScopeMap[t]
		typeName, found := c.TypeMap[t]
		if !found {
			typeName = t
		}
		ret += fmt.Sprintf("#### %s\n\n%s\n", typeName, msg)
	}
	return ret
}

// DefaultFormatFunc is used to format a commit message
func DefaultFormatFunc(c *repository.Commit) string {
	if c.Scope != "" {
		return fmt.Sprintf("* %s [%s] (%s) \n", c.Subject, c.Scope, TrimSHA(c.Hash))
	}
	return fmt.Sprintf("* %s (%s) \n", c.Subject, TrimSHA(c.Hash))
}

// URLFormatFunc is used to format a commit message
func URLFormatFunc(url string) FormatFunc {
	return func(c *repository.Commit) string {
		if c.Scope != "" {
			ticketURL := strings.Replace(url, "{SCOPE}", c.Scope, -1)
			return fmt.Sprintf("* %s [%s](%s) (%s) \n", c.Subject, c.Scope, ticketURL, TrimSHA(c.Hash))
		}
		return DefaultFormatFunc(c)
	}
}

// TrimSHA returns only the leading 8 characters of a commit hash
func TrimSHA(sha string) string {
	if len(sha) < 9 {
		return sha
	}
	return sha[:8]
}

func getSortedKeys(m *map[string]string) []string {
	keys := make([]string, len(*m))
	i := 0
	for k := range *m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
