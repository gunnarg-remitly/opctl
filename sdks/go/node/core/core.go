// Package core defines the core interface for an opspec node
package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/dgraph-io/badger/v2"
	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o fakes/core.go . Core
type Core interface {
	AddAuth(
		req model.AddAuthReq,
	) error
	GetEventStream(
		ctx context.Context,
		req *model.GetEventStreamReq,
	) (
		<-chan model.Event,
		<-chan error,
	)

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
) (Core, error) {
	eventDbPath := path.Join(dataDirPath, "dcg", "events")
	if err := os.MkdirAll(eventDbPath, 0700); nil != err {
		return nil, err
	}

	// per badger README.MD#FAQ "maximizes throughput"
	runtime.GOMAXPROCS(128)

	db, err := badger.Open(
		badger.DefaultOptions(
			eventDbPath,
		).WithLogger(nil),
	)
	if err != nil {
		return nil, err
	}

	pubSub := pubsub.New(db)

	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	stateStore := newStateStore(
		ctx,
		db,
		pubSub,
	)

	caller := newCaller(
		newContainerCaller(
			containerRuntime,
			pubSub,
			stateStore,
		),
		dataDirPath,
		stateStore,
		pubSub,
	)

	return _core{
		caller:           caller,
		containerRuntime: containerRuntime,
		dataCachePath:    filepath.Join(dataDirPath, "ops"),
		opCaller: newOpCaller(
			stateStore,
			caller,
			dataDirPath,
		),
		pubSub:              pubSub,
		stateStore:          stateStore,
		uniqueStringFactory: uniqueStringFactory,
	}, nil
}

type _core struct {
	caller              caller
	containerRuntime    containerruntime.ContainerRuntime
	dataCachePath       string
	opCaller            opCaller
	pubSub              pubsub.PubSub
	stateStore          stateStore
	uniqueStringFactory uniquestring.UniqueStringFactory
}
