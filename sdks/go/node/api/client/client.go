// Package client implements a client for the opspec node api
package client

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	iwebsocket "github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/gorilla/websocket"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/k8s"
)

//counterfeiter:generate -o fakes/client.go . Client
type Client interface {
	// AddAuth adds auth
	AddAuth(
		ctx context.Context,
		req model.AddAuthReq,
	) error

	GetEventStream(
		ctx context.Context,
		req *model.GetEventStreamReq,
	) (
		stream <-chan model.Event,
		err error,
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

	KillOp(
		ctx context.Context,
		req model.KillOpReq,
	) (
		err error,
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

	// Liveness ensures liveness of the node
	Liveness(
		ctx context.Context,
	) error

	StartOp(
		ctx context.Context,
		req model.StartOpReq,
	) (
		opID string,
		err error,
	)
}

type Opts struct {
	// RetryLogHook will be executed anytime a request is retried
	RetryLogHook     func(err error)
	ContainerRuntime string
	DataDirPath      string
}

// New returns a new client
// nil opts will be ignored
func New(
	opts *Opts,
) Client {
	var containerRuntime containerruntime.ContainerRuntime
	var err error
	if "k8s" == opts.ContainerRuntime {
		containerRuntime, err = k8s.New()
		if nil != err {
			panic(err)
		}
	} else {
		containerRuntime, err = docker.New()
		if nil != err {
			panic(err)
		}
	}
	c := core.New(containerRuntime, opts.DataDirPath)

	return &client{
		core:     c,
		wsDialer: websocket.DefaultDialer,
	}
}

type client struct {
	core     core.Core
	wsDialer iwebsocket.Dialer
}
