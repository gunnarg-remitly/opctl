package core

import (
	"context"
	"fmt"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
)

// Eventser exposes the "events" command
type Eventser interface {
	Events(ctx context.Context)
}

// newEventser returns an initialized "events" command
func newEventser(
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
	core core.Core,
) Eventser {
	return _eventser{
		cliExiter: cliExiter,
		cliOutput: cliOutput,
		core:      core,
	}
}

type _eventser struct {
	cliExiter cliexiter.CliExiter
	cliOutput clioutput.CliOutput
	core      core.Core
}

func (ivkr _eventser) Events(
	ctx context.Context,
) {
	eventChannel, errChan := ivkr.core.GetEventStream(
		ctx,
		&model.GetEventStreamReq{},
	)
	go func() {
		for {
			err := <-errChan
			if err != nil {
				fmt.Printf("error received: %v\n", err)
			}
		}
	}()

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
