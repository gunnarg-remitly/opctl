package auth

import (
	"context"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
)

// Adder exposes the "auth add" sub command
type Adder interface {
	Add(
		ctx context.Context,
		resources string,
		username string,
		password string,
	)
}

// newAdder returns an initialized "auth add" sub command
func newAdder(
	cliExiter cliexiter.CliExiter,
	core core.Core,
) Adder {
	return _adder{
		cliExiter: cliExiter,
		core:      core,
	}
}

type _adder struct {
	cliExiter cliexiter.CliExiter
	core      core.Core
}

func (ivkr _adder) Add(
	ctx context.Context,
	resources string,
	username string,
	password string,
) {
	err := ivkr.core.AddAuth(
		model.AddAuthReq{
			Resources: resources,
			Creds: model.Creds{
				Username: username,
				Password: password,
			},
		},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
