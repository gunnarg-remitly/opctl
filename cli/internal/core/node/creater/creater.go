package creater

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/cli/internal/runtime"
	"github.com/opctl/opctl/sdks/go/node/core"
)

// Creater exposes the "node create" sub command
//counterfeiter:generate -o fakes/creater.go . Creater
type Creater interface {
	Create(
		opts nodeprovider.NodeOpts,
	) error
}

// New returns an initialized "node create" command
func New() Creater {
	return _creater{}
}

type _creater struct{}

func (ivkr _creater) Create(
	opts nodeprovider.NodeOpts,
) error {
	dataDir, err := datadir.New(opts.DataDir)
	if nil != err {
		return err
	}

	if err := dataDir.InitAndLock(); nil != err {
		return err
	}

	containerRuntime, err := runtime.New(opts.ContainerRuntime)

	err = newHTTPListener(
		core.New(
			containerRuntime,
			dataDir.Path(),
		),
	).
		Listen(
			context.Background(),
			opts.ListenAddress,
		)

	if nil != err {
		return err
	}

	return nil
}
