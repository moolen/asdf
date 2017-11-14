package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/moolen/asdf/repository"
)

func readVersionFile(path string) (*semver.Version, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return readVersion(file)
}

func readVersion(rd io.Reader) (*semver.Version, error) {
	content, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	versionString := strings.TrimRight(string(content), "\n")
	version := semver.MustParse(versionString)
	return version, nil
}

func nextReleaseByChange(latest *semver.Version, change repository.Change) semver.Version {
	switch change {
	case repository.MajorChange:
		log.Printf("increment major")
		return latest.IncMajor()
	case repository.MinorChange:
		log.Printf("increment minor")
		return latest.IncMinor()
	}
	return latest.IncPatch()
}
