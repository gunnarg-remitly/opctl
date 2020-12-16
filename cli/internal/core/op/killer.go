package op

import (
	"context"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// Killer exposes the "op kill" sub command
type Killer interface {
	Kill(
		ctx context.Context,
		opID string,
	)
}

// newKiller returns an initialized "op kill" sub command
func newKiller(
	cliExiter cliexiter.CliExiter,
	api client.Client,
) Killer {
	return _killer{
		cliExiter: cliExiter,
		api:       api,
	}
}

type _killer struct {
	cliExiter cliexiter.CliExiter
	api       client.Client
}

func (ivkr _killer) Kill(
	ctx context.Context,
	opID string,
) {
	err := ivkr.api.KillOp(
		ctx,
		model.KillOpReq{
			OpID:       opID,
			RootCallID: opID,
		},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
