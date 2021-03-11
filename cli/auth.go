package main

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
)

// auth implements "auth" command
func auth(
	ctx context.Context,
	node core.Core,
	addAuthReq model.AddAuthReq,
) error {
	return node.AddAuth(
		ctx,
		addAuthReq,
	)
}
