package core

import (
	"context"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c _core) GetData(
	ctx context.Context,
	req model.GetDataReq,
) (
	model.ReadSeekCloser,
	error,
) {
	if req.PkgRef == "" || req.ContentPath == "" {
		return nil, fmt.Errorf("not found: %s%s", req.PkgRef, req.ContentPath)
	}

	dataHandle, err := c.ResolveData(ctx, req.PkgRef, req.PullCreds)
	if err != nil {
		return nil, err
	}

	// this might not be right
	return dataHandle.GetContent(ctx, req.ContentPath)
}
