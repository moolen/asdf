package main

import (
	"io/ioutil"
	"path/filepath"
	"io"
	"fmt"
	"net/http"
	"github.com/Masterminds/semver"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/google/go-github/github"
	"context"
	"os"
	"golang.org/x/oauth2"
	"runtime"
)

func syncAction(c *cli.Context) error{
	var explicitVersion *semver.Version
	release := c.String(flagRelease)
	if release != "" {
		var err error
		log.Debugf("starting sync for explicit release version %s", release)
		explicitVersion, err = semver.NewVersion(release)
		if err != nil {
			log.Errorf("not a semver: %s", release)
			return cli.NewExitError(err, 1)
		}
	}
	accessToken := c.String(flagAccessToken)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	releases, _, err := client.Repositories.ListReleases(ctx, "figome", "semantic-changelog", nil)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if len(releases) == 0 {
		return cli.NewExitError("no release found", 2)
	}
	maxVersion := semver.MustParse("0.0.0")
	var myRelease *github.RepositoryRelease 
	for _, release := range releases {
		log.Infof("found release %s / tag %s", release.GetName(), release.GetTagName())

		thisVersion, err := semver.NewVersion(release.GetTagName())
		if err != nil {
			log.Infof("skipping release, not a semver: %s", release.GetTagName())
			continue
		}
		if explicitVersion == nil {
			if thisVersion.GreaterThan(maxVersion) {
				maxVersion = thisVersion
				myRelease = release
			}
			continue
		}
		if explicitVersion.Equal(thisVersion) {
			myRelease = release
		}
	}

	if myRelease == nil {
		return cli.NewExitError(fmt.Errorf("no suitable release found"), 3)
	}

	log.Infof("using release tag %s", myRelease.GetTagName())
	var downloadURL string
	for _, asset := range myRelease.Assets {
		if strings.Contains(asset.GetName(), runtime.GOOS) {
			log.Debugf("found download url for %s: %s", asset.GetName(), asset.GetBrowserDownloadURL())
			downloadURL = asset.GetBrowserDownloadURL()
		}
	}
	if downloadURL == "" {
		return cli.NewExitError(fmt.Errorf("could not find suitable download url"), 4)
	}
	log.Infof("found download url release tag %s", myRelease.GetTagName())
	ex, err := os.Executable()
	if err != nil {
		return cli.NewExitError(fmt.Errorf("could not get executable: %s", err), 8)
	}
	tmpPath := filepath.Dir(ex)
	tmpFile, err := ioutil.TempFile(tmpPath, "")
	if err != nil {
		return cli.NewExitError(fmt.Errorf("could not create temp file for downloading"), 5)
	}
	resp, err := http.Get(downloadURL)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("could fetch download: %s", err), 6)
	}
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("could not copy http response body: %s", err), 7)
	}
	resp.Body.Close()
	tmpFile.Close()
	
	log.Infof("renaming %s to %s", tmpFile.Name(), ex)
	err = os.Rename(tmpFile.Name(), ex)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("could not rename: %s", err), 9)
	}
	log.Infof("setting chmod for: %s", ex)
	err = os.Chmod(ex, 0777)
	if err != nil {
		return cli.NewExitError(err, 10)
	}
	log.Infof("new binary available: %s", ex)
	return nil
}

func syncFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  flagRelease,
			Value: "",
			Usage: "use to fetch a specific version",
		},
		cli.StringFlag{
			Name: flagAccessToken,
			Value: "",
			Usage: "your personal github access token",
			EnvVar: "GITHUB_ACCESS_TOKEN",
		},
	}
}