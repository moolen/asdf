package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

const (
	DefaultVersionFile   = "VERSION"
	DefaultChangelogFile = "CHANGELOG.md"
)

// Config holds all the configuration
type Config struct {
	VersionFile   string          `json:"version_file,omitempty"`
	ChangelogFile string          `json:"changelog_file,omitempty"`
	Repository    string          `json:"repository"`
	BranchSuffix  BranchSuffix    `json:"branch_suffix"`
	Types         TypeConstraints `json:"types,omitempty"`
	TicketURL     string          `json:"ticket_url,omitempty"`
}

// BranchSuffix contains a mapping between
// branch-names and a prerelease scheme
type BranchSuffix map[string]string

// FromJSON returns a config from a io.Reader that contains json-encoded data
func FromJSON(reader io.Reader) (*Config, error) {
	config := &Config{
		VersionFile:   DefaultVersionFile,
		ChangelogFile: DefaultChangelogFile,
	}
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	if config.Types == nil {
		config.Types = defaultTypes()
	}
	return config, nil
}

// FromFile tries to read a file containing JSON
func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return FromJSON(file)
}

func defaultTypes() TypeConstraints {
	return TypeConstraints{
		&TypeConstraint{
			Key:   "feat",
			Label: "Feature",
			Minor: true,
		},
		&TypeConstraint{
			Key:   "breaking",
			Label: "Breaking Changes",
			Major: true,
		},
		&TypeConstraint{
			Key:   "fix",
			Label: "Bug Fixes",
		},
		&TypeConstraint{
			Key:   "perf",
			Label: "Performance Improvements",
		},
		&TypeConstraint{
			Key:   "revert",
			Label: "Reverts",
		},
		&TypeConstraint{
			Key:   "docs",
			Label: "Documentation",
		},
		&TypeConstraint{
			Key:   "refactor",
			Label: "Code Refactoring",
		},
		&TypeConstraint{
			Key:   "test",
			Label: "Tests",
		},
		&TypeConstraint{
			Key:   "chore",
			Label: "Chores",
		},
	}
}
