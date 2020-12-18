package core

import (
	"time"

	"github.com/opctl/opctl/sdks/go/model"
)

func (this _core) AddAuth(
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
			Timestamp: time.Now().UTC(),
		},
	)
	return nil
}
