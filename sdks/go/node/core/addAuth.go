package core

import (
	"github.com/opctl/opctl/sdks/go/model"
)

func (this _core) AddAuth(
	req model.AddAuthReq,
) error {
	return this.stateStore.AddAuth(
		model.AuthAdded{
			Auth: model.Auth{
				Creds:     req.Creds,
				Resources: req.Resources,
			},
		},
	)
}
