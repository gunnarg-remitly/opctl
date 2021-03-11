package main

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func selfUpdate() (string, error) {
	v := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(v, "opctl/opctl")
	if nil != err {
		return "", err
	}

	if latest.Version.Equals(v) {
		return "No update available, already at the latest version!", nil
	}

	return fmt.Sprintf("Updated to new version: %s!", latest.Version), err
}
