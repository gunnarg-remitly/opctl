package client

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c client) KillOp(
	ctx context.Context,
	req model.KillOpReq,
) error {
	return c.core.KillOp(req)
}
