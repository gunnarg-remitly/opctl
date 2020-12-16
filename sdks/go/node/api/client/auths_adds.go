package client

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c client) AddAuth(
	ctx context.Context,
	req model.AddAuthReq,
) error {
	c.core.AddAuth(req)
	return nil
}
