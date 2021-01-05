// Package core defines the core interface for an opspec node
package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"path/filepath"
	"runtime"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
)

//counterfeiter:generate -o fakes/core.go . Core
type Core interface {
	AddAuth(
		req model.AddAuthReq,
	) error

	StartOp(
		ctx context.Context,
		req model.StartOpReq,
	) (
		callID string,
		err error,
	)

	// Resolve attempts to resolve an op via local filesystem or git
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

	// ListDescendants lists file system entries
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	ListDescendants(
		ctx context.Context,
		req model.ListDescendantsReq,
	) (
		[]*model.DirEntry,
		error,
	)

	// GetData gets data
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	GetData(
		ctx context.Context,
		req model.GetDataReq,
	) (
		model.ReadSeekCloser,
		error,
	)
}

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

	return _core{
		caller:              caller,
		dataCachePath:       filepath.Join(dataDirPath, "ops"),
		stateStore:          stateStore,
		uniqueStringFactory: uniqueStringFactory,
	}, nil
}

type _core struct {
	caller              caller
	dataCachePath       string
	stateStore          stateStore
	uniqueStringFactory uniquestring.UniqueStringFactory
}
