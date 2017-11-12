package config

import (
	"reflect"
	"testing"

	"github.com/moolen/asdf/repository"
)

func TestTypeConstraint(t *testing.T) {

	table := []struct {
		t      *TypeConstraint
		change Change
	}{
		{
			t: &TypeConstraint{
				Major: true,
			},
			change: ChangeMajor,
		},
		{
			t: &TypeConstraint{
				Minor: true,
			},
			change: ChangeMinor,
		},
		{
			t: &TypeConstraint{
				Patch: true,
			},
			change: ChangePatch,
		},
		{
			t:      &TypeConstraint{},
			change: ChangePatch,
		},
	}

	for i, row := range table {
		change := row.t.Max()
		if change != row.change {
			t.Fatalf("[%d] wrong change: expected %d, got %d", i, row.change, change)
		}
	}
}

func TestTypeConstraintsKeyChangeMap(t *testing.T) {

	table := []struct {
		t TypeConstraints
		m map[string]Change
	}{
		{
			t: TypeConstraints{
				&TypeConstraint{
					Key:   "mykey",
					Major: true,
				},
			},
			m: map[string]Change{
				"mykey": ChangeMajor,
			},
		},
		{
			t: TypeConstraints{
				&TypeConstraint{
					Key:   "mykey",
					Minor: true,
				},
			},
			m: map[string]Change{
				"mykey": ChangeMinor,
			},
		},
		{
			t: TypeConstraints{
				&TypeConstraint{
					Key:   "mykey",
					Patch: true,
				},
			},
			m: map[string]Change{
				"mykey": ChangePatch,
			},
		},
	}

	for i, row := range table {
		m := row.t.KeyChangeMap()
		if !reflect.DeepEqual(m, row.m) {
			t.Fatalf("[%d] table: expected \n%#v, got \n%#v", i, row.m, m)
		}
	}
}

func TestTypeConstraintsKeyLabelMap(t *testing.T) {

	table := []struct {
		t TypeConstraints
		m map[string]string
	}{
		{
			t: TypeConstraints{
				&TypeConstraint{
					Key:   "mykey",
					Label: "bart",
				},
			},
			m: map[string]string{
				"mykey": "bart",
			},
		},
		{
			t: TypeConstraints{
				&TypeConstraint{
					Key:   "fart",
					Label: "foo",
					Minor: true,
				},
			},
			m: map[string]string{
				"fart": "foo",
			},
		},
	}

	for i, row := range table {
		m := row.t.KeyLabelMap()
		if !reflect.DeepEqual(m, row.m) {
			t.Fatalf("[%d] table: expected \n%#v, got \n%#v", i, row.m, m)
		}
	}
}

func TestTypeConstraintsMax(t *testing.T) {
	table := []struct {
		t       TypeConstraints
		commits []*repository.Commit
		change  Change
	}{
		{
			t: TypeConstraints{
				&TypeConstraint{
					Key:   "feat",
					Minor: true,
				},
			},
			commits: []*repository.Commit{
				&repository.Commit{
					Type: "feat",
				},
			},
			change: ChangeMinor,
		},
		{
			t: TypeConstraints{},
			commits: []*repository.Commit{
				&repository.Commit{
					Type: "breaking",
				},
			},
			change: ChangePatch,
		},
		{
			t: TypeConstraints{
				&TypeConstraint{
					Key:   "docs",
					Major: true,
				},
			},
			commits: []*repository.Commit{
				&repository.Commit{
					Type: "feat",
				},
				&repository.Commit{
					Type: "docs",
				},
			},
			change: ChangeMajor,
		},
	}

	for i, row := range table {
		max := row.t.Max(row.commits)
		if max != row.change {
			t.Fatalf("[%d] table: expected %#v, got %#v", i, row.change, max)
		}
	}
}
