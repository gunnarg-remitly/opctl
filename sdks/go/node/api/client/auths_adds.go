package client

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c client) AddAuth(
	ctx context.Context,
	req model.AddAuthReq,
) error {
	return c.core.AddAuth(req)
}
