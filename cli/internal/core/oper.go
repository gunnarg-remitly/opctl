package core

import (
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/core/op"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// Oper exposes the "op" sub command
type Oper interface {
	Op() op.Op
}

// newOper returns an initialized "op" sub command
func newOper(
	cliExiter cliexiter.CliExiter,
	dataResolver dataresolver.DataResolver,
	api client.Client,
) Oper {
	return _oper{
		op: op.New(
			cliExiter,
			dataResolver,
			api,
		),
	}
}

type _oper struct {
	op op.Op
}

func (ivkr _oper) Op() op.Op {
	return ivkr.op
}
