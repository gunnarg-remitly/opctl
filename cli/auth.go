package main

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

// auth implements "auth" command
func auth(
	ctx context.Context,
	addAuthReq model.AddAuthReq,
) error {
	return node.AddAuth(
		ctx,
		addAuthReq,
	)
}
