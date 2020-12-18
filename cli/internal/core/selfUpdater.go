package core

import (
	"fmt"

	"github.com/opctl/opctl/cli/internal/updater"
)

// SelfUpdater exposes the "self-update" command
type SelfUpdater interface {
	SelfUpdate(
		releaseChannel string,
	) (string, error)
}

// newSelfUpdater returns an initialized "self-update" command
func newSelfUpdater() SelfUpdater {
	return _selfUpdateInvoker{
		updater: updater.New(),
	}
}

type _selfUpdateInvoker struct {
	updater updater.Updater
}

func (ivkr _selfUpdateInvoker) SelfUpdate(
	releaseChannel string,
) (string, error) {
	if releaseChannel != "alpha" && releaseChannel != "beta" && releaseChannel != "stable" {
		return "", fmt.Errorf(
			"%v is not an available release channel. "+
				"Available release channels are 'alpha', 'beta', and 'stable'.", releaseChannel)
	}

	update, err := ivkr.updater.GetUpdateIfExists(releaseChannel)
	if nil != err {
		return "", err
	} else if nil == update {
		return "No update available, already at the latest version!", err
	}

	err = ivkr.updater.ApplyUpdate(update)
	if nil != err {
		return "", err
	}

	// @TODO start node maintaining previous user
	return fmt.Sprintf("Updated to new version: %s!", update.Version), nil
}
