package client

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c client) KillOp(
	ctx context.Context,
	req model.KillOpReq,
) error {
	c.core.KillOp(req)
	return nil
}
