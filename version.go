package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Masterminds/semver"
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
