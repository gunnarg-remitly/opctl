package core

import (
	"context"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// Eventser exposes the "events" command
type Eventser interface {
	Events(ctx context.Context)
}

// newEventser returns an initialized "events" command
func newEventser(
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
	api client.Client,
) Eventser {
	return _eventser{
		cliExiter: cliExiter,
		cliOutput: cliOutput,
		api:       api,
	}
}

type _eventser struct {
	cliExiter cliexiter.CliExiter
	cliOutput clioutput.CliOutput
	api       client.Client
}

func (ivkr _eventser) Events(
	ctx context.Context,
) {
	eventChannel, err := ivkr.api.GetEventStream(
		ctx,
		&model.GetEventStreamReq{},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for {
		event, isEventChannelOpen := <-eventChannel
		if !isEventChannelOpen {
			ivkr.cliExiter.Exit(
				cliexiter.ExitReq{
					Message: "Connection to event stream lost",
					Code:    1,
				},
			)
			return // support fake exiter
		}

		ivkr.cliOutput.Event(&event)
	}
}
