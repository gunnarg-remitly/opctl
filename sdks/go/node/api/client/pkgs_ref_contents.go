package client

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c client) ListDescendants(
	ctx context.Context,
	req model.ListDescendantsReq,
) (
	[]*model.DirEntry,
	error,
) {
	if req.PkgRef == "" {
		return []*model.DirEntry{}, nil
	}

	dataHandle, err := c.core.ResolveData(ctx, req.PkgRef, req.PullCreds)
	if err != nil {
		return nil, err
	}

	// this might not be right
	return dataHandle.ListDescendants(ctx)
}
