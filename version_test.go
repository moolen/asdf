package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestReadVersionFile(t *testing.T) {
	_, err := readVersionFile("")
	if err == nil {
		t.Fail()
	}
	repo := createRepository()
	versionFile := path.Join(repo, "VERSION")
	version, err := readVersionFile(versionFile)
	if err != nil {
		t.Fail()
	}
	if version.String() != "1.0.0" {
		t.Fail()
	}
	ioutil.WriteFile(versionFile, []byte("2.1.31"), os.ModePerm)
	version, err = readVersionFile(versionFile)
	if err != nil {
		t.Fail()
	}
	if version.String() != "2.1.31" {
		t.Fail()
	}
}
