package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
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

func calcReleaseVersion(latestCommit *repository.Commit, branch string, latest *semver.Version, branchSuffix map[string]string, change repository.Change) (*semver.Version, error) {
	var next semver.Version
	var err error
	log.Printf("latest release: %s", latest)
	// if we're on master we'll either:
	// - remove the prerelease suffix
	// - or increment to next version
	if branch == "master" {
		if latest.Prerelease() == "" {
			log.Printf("on branch master without prerelease, calculating next change: %d", change)
			next = nextReleaseByChange(latest, change)
		} else {
			next, err = latest.SetPrerelease("")
			log.Printf("on branch master with prerelease, removing prerelease")
			if err != nil {
				return nil, err
			}
		}
		return &next, nil
	}

	// if this branch has a special mapping
	// we apply the mapping withouth incrementing
	for branchRE, suffix := range branchSuffix {
		re := regexp.MustCompile(branchRE)
		if re.MatchString(branch) {
			nextSuffix := nextPrereleaseSuffix(latest, latestCommit, suffix)
			next, err = latest.SetPrerelease(nextSuffix)
		}
	}
	return &next, err
}

func nextPrereleaseSuffix(latest *semver.Version, commit *repository.Commit, suffix string) string {
	var nextSuffix string
	hash := commit.Hash
	if len(hash) > 8 {
		hash = hash[:8]
	}
	prerelease := latest.Prerelease()
	if strings.Contains(suffix, ReleaseToken) {
		stripped := strings.Replace(suffix, ReleaseToken, "", 1)
		strippedPrerelease := strings.Replace(prerelease, stripped, "", 1)
		releaseNumber, _ := strconv.ParseInt(strippedPrerelease, 10, 64)
		nextReleaseNum := fmt.Sprintf("%d", releaseNumber+1)
		nextSuffix = strings.Replace(suffix, ReleaseToken, nextReleaseNum, 1)
	} else if strings.Contains(suffix, CommitToken) {
		nextSuffix = strings.Replace(suffix, "{COMMIT_SHA}", hash, 1)
	}
	return nextSuffix
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
