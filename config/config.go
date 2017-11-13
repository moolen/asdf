package config

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

const (
	DefaultVersionFile   = "VERSION"
	DefaultChangelogFile = "CHANGELOG.md"
)

// ErrNoConfigFile happens if the config file does not exist
var ErrNoConfigFile = errors.New("Config file does not exist")

// Config holds all the configuration
type Config struct {
	VersionFile   string            `json:"version_file,omitempty"`
	ChangelogFile string            `json:"changelog_file,omitempty"`
	Repository    string            `json:"repository"`
	BranchSuffix  BranchSuffix      `json:"branch_suffix"`
	Types         map[string]string `json:"types,omitempty"`
	TicketURL     string            `json:"ticket_url,omitempty"`
}

// BranchSuffix contains a mapping between
// branch-names and a prerelease scheme
type BranchSuffix map[string]string

// DefaultTypeMap contains a mapping of types to groups
// which are used to render the changelog
var DefaultTypeMap = map[string]string{
	"feat":     "Feature",
	"breaking": "Breaking Changes",
	"fix":      "Bug Fixes",
	"perf":     "Performance Improvements",
	"revert":   "Reverted",
	"docs":     "Documentation",
	"refactor": "Code Refactoring",
	"test":     "Tests",
	"chore":    "Chores",
}

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
		config.Types = DefaultTypeMap
	}
	return config, nil
}

// FromFile tries to read a file containing JSON
func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, ErrNoConfigFile
	}
	return FromJSON(file)
}
