package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
)

// New returns a data provider which sources pkgs from a node
// A node now represents a local installation of opctl, where pkgs can be
// installed into opctl's data directory.
func New(
	core core.Core,
	pullCreds *model.Creds,
) model.DataProvider {
	return _node{
		core:      core,
		pullCreds: pullCreds,
	}
}

type _node struct {
	core      core.Core
	pullCreds *model.Creds
}

func (np _node) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	// ensure resolvable by listing contents w/out err
	if _, err := np.core.ListDescendants(
		ctx,
		model.ListDescendantsReq{
			PkgRef:    dataRef,
			PullCreds: np.pullCreds,
		},
	); nil != err {
		return nil, err
	}

	return newHandle(np.core, dataRef, np.pullCreds), nil
}
