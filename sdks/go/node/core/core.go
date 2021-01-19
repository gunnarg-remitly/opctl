// Package core defines the core interface for an opspec node
package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"path/filepath"
	"runtime"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
)

// New returns a new LocalCore initialized with the given options
func New(
	ctx context.Context,
	containerRuntime containerruntime.ContainerRuntime,
	dataDirPath string,
	eventChannel chan model.Event,
) (Core, error) {
	// per badger README.MD#FAQ "maximizes throughput"
	runtime.GOMAXPROCS(128)

	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	stateStore, err := newStateStore(
		ctx,
		dataDirPath,
	)
	if err != nil {
		return nil, err
	}

	caller := newCaller(
		newContainerCaller(
			containerRuntime,
			eventChannel,
			stateStore,
		),
		dataDirPath,
		eventChannel,
	)

	return core{
		caller:              caller,
		dataCachePath:       filepath.Join(dataDirPath, "ops"),
		stateStore:          stateStore,
		uniqueStringFactory: uniqueStringFactory,
	}, nil
}

// core is an OpNode that supports running ops directly on the host
type core struct {
	caller              caller
	dataCachePath       string
	stateStore          stateStore
	uniqueStringFactory uniquestring.UniqueStringFactory
}

//counterfeiter:generate -o fakes/core.go . Core

// Core is an OpNode that supports running ops directly on the current machine
type Core interface {
	node.OpNode

	// Resolve attempts to resolve data via local filesystem or git
	// nil pullCreds will be ignored
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	ResolveData(
		ctx context.Context,
		dataRef string,
		pullCreds *model.Creds,
	) (
		model.DataHandle,
		error,
	)
}
