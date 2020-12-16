package client

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

// StartOp starts an op & returns its root op id (ROId)
func (c client) StartOp(
	ctx context.Context,
	req model.StartOpReq,
) (string, error) {

	callID, err := c.core.StartOp(ctx, req)
	if nil != err {
		return "", err
	}

	return string(callID), nil
}
