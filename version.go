package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/Masterminds/semver"
	"github.com/figome/semantic-changelog/repository"
)

func readVersionFile(path string) (*semver.Version, error) {
	log.Debugf("trying to open %s", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err // todo: more descriptive error message
	}
	return readVersion(file)
}

func readVersion(rd io.Reader) (*semver.Version, error) {
	content, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	versionString := strings.TrimRight(string(content), "\n")
	log.Debugf("file content %s", versionString)
	version, err := semver.NewVersion(versionString)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func nextReleaseByChange(latest *semver.Version, change repository.Change) semver.Version {
	switch change {
	case repository.MajorChange:
		log.Debugf("increment major")
		return latest.IncMajor()
	case repository.MinorChange:
		log.Debugf("increment minor")
		return latest.IncMinor()
	}
	return latest.IncPatch()
}
