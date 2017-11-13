package main

import (
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/moolen/asdf/repository"
)

func TestCalcNextVersion(t *testing.T) {
	table := []struct {
		commit      *repository.Commit
		branch      string
		version     *semver.Version
		nextVersion *semver.Version
		suffixMap   map[string]string
		change      repository.Change
		err         error
	}{
		{
			commit:      &repository.Commit{},
			branch:      "master",
			version:     semver.MustParse("1.2.3"),
			nextVersion: semver.MustParse("1.2.4"),
			change:      repository.PatchChange,
		},
		{
			commit:      &repository.Commit{},
			branch:      "master",
			version:     semver.MustParse("1.2.3-rc400"),
			nextVersion: semver.MustParse("1.2.3"),
			change:      repository.PatchChange,
		},
		{
			commit: &repository.Commit{
				Hash: "1234",
			},
			branch:      "devrelease",
			version:     semver.MustParse("1.2.3"),
			nextVersion: semver.MustParse("1.2.3-dev1234"),
			suffixMap: map[string]string{
				"devrelease": "dev{COMMIT_SHA}",
			},
			change: repository.MinorChange,
		},
		{
			commit:      &repository.Commit{},
			branch:      "release",
			version:     semver.MustParse("1.2.3-rc1"),
			nextVersion: semver.MustParse("1.2.3-rc2"),
			suffixMap: map[string]string{
				"release": "rc{RELEASE_NUMBER}",
			},
			change: repository.PatchChange,
		},
		{
			commit:      &repository.Commit{},
			branch:      "beta",
			version:     semver.MustParse("2.0.0-beta.1"),
			nextVersion: semver.MustParse("2.0.0-beta.2"),
			suffixMap: map[string]string{
				"beta": "beta.{RELEASE_NUMBER}",
			},
			change: repository.PatchChange,
		},
	}

	for i, row := range table {
		next, err := calcReleaseVersion(row.commit, row.branch, row.version, row.suffixMap, row.change)
		if err != row.err {
			t.Fatalf("[%d] expected %s\ngot %s", i, row.err, err)
		}
		if !reflect.DeepEqual(row.nextVersion, next) {
			t.Fatalf("[%d] expected %s\ngot %s", i, row.nextVersion, next)
		}
	}
}
