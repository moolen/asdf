package config

import "github.com/moolen/asdf/repository"

// Change represents an semantic change
// Change is internaly an int so we can easily sort a list of changes:
// e.g. MajorChange > PatchChange
type Change int

const (
	// ChangePatch for a patch release
	ChangePatch Change = 0
	// ChangeMinor for a minor release
	ChangeMinor = 1
	// ChangeMajor for a major release
	ChangeMajor = 2
)

// TypeConstraint is a (optional) part of the asdf config
// they define which commit "type" triggers what kind of change
type TypeConstraint struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Major bool   `json:"major,omitempty"`
	Minor bool   `json:"minor,omitempty"`
	Patch bool   `json:"patch,omitempty"`
}

// TypeConstraints provides auxiliary methods
// to determine the severity of a change in a list of commits
type TypeConstraints []*TypeConstraint

// Max returns the biggest change from a slice of commits
// if no commits are present it will  return a PatchChange
func (t TypeConstraints) Max(commits []*repository.Commit) Change {
	max := ChangePatch
	label := t.KeyChangeMap()
	for _, commit := range commits {
		if label[commit.Type] > max {
			max = label[commit.Type]
		}
	}
	return max
}

// KeyLabelMap returns a map[key]label
func (t TypeConstraints) KeyLabelMap() map[string]string {
	m := make(map[string]string)
	for _, c := range t {
		m[c.Key] = c.Label
	}
	return m
}

// KeyChangeMap returns a map[key]change
func (t TypeConstraints) KeyChangeMap() map[string]Change {
	m := make(map[string]Change)
	for _, c := range t {
		m[c.Key] = c.Max()
	}
	return m
}

// Max gives us the severity of the change of this kind
func (t TypeConstraint) Max() Change {
	switch {
	case t.Major:
		return ChangeMajor
	case t.Minor:
		return ChangeMinor
	}
	return ChangePatch
}
