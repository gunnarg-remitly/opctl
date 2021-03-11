package core

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

func (this core) AddAuth(
	ctx context.Context,
	req model.AddAuthReq,
) error {
	this.pubSub.Publish(
		model.Event{
			AuthAdded: &model.AuthAdded{
				Auth: model.Auth{
					Creds:     req.Creds,
					Resources: req.Resources,
				},
			},
		},
	)
	return nil
}
