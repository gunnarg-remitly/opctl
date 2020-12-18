package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	cliModel "github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// Runer exposes the "run" command
type Runer interface {
	Run(
		ctx context.Context,
		opRef string,
		opts *cliModel.RunOpts,
	)
}

// newRuner returns an initialized "run" command
func newRuner(
	cliColorer clicolorer.CliColorer,
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	dataResolver dataresolver.DataResolver,
	core core.Core,
) Runer {
	return _runer{
		cliColorer:        cliColorer,
		cliExiter:         cliExiter,
		cliOutput:         cliOutput,
		cliParamSatisfier: cliParamSatisfier,
		dataResolver:      dataResolver,
		core:              core,
	}
}

type _runer struct {
	dataResolver      dataresolver.DataResolver
	cliColorer        clicolorer.CliColorer
	cliExiter         cliexiter.CliExiter
	cliOutput         clioutput.CliOutput
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier
	core              core.Core
}

func (ivkr _runer) Run(
	ctx context.Context,
	opRef string,
	opts *cliModel.RunOpts,
) {
	startTime := time.Now().UTC()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opHandle := ivkr.dataResolver.Resolve(
		ctx,
		opRef,
		nil,
	)

	opFileReader, err := opHandle.GetContent(
		ctx,
		opfile.FileName,
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	opFileBytes, err := ioutil.ReadAll(opFileReader)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	opFile, err := opfile.Unmarshal(
		opFileBytes,
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	ymlFileInputSrc, err := ivkr.cliParamSatisfier.NewYMLFileInputSrc(opts.ArgFile)
	if nil != err {
		err = fmt.Errorf("unable to load arg file at '%v'; error was: %v", opts.ArgFile, err.Error())
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	argsMap := ivkr.cliParamSatisfier.Satisfy(
		cliparamsatisfier.NewInputSourcer(
			ivkr.cliParamSatisfier.NewSliceInputSrc(opts.Args, "="),
			ymlFileInputSrc,
			ivkr.cliParamSatisfier.NewEnvVarInputSrc(),
			ivkr.cliParamSatisfier.NewParamDefaultInputSrc(opFile.Inputs),
			ivkr.cliParamSatisfier.NewCliPromptInputSrc(opFile.Inputs),
		),
		opFile.Inputs,
	)

	// init signal channel
	aSigIntWasReceivedAlready := false
	sigIntChannel := make(chan os.Signal, 1)
	defer close(sigIntChannel)
	signal.Notify(
		sigIntChannel,
		syscall.SIGINT,
	)

	sigTermChannel := make(chan os.Signal, 1)
	defer close(sigTermChannel)
	signal.Notify(
		sigTermChannel,
		syscall.SIGTERM,
	)

	// start op
	rootCallID, err := ivkr.core.StartOp(
		ctx,
		model.StartOpReq{
			Args: argsMap,
			Op: model.StartOpReqOp{
				Ref: opHandle.Ref(),
			},
		},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	// start event loop
	eventChannel, errChan := ivkr.core.GetEventStream(
		ctx,
		&model.GetEventStreamReq{
			Filter: model.EventFilter{
				Roots: []string{rootCallID},
				Since: &startTime,
			},
		},
	)

	for {
		select {
		case err := <-errChan:
			fmt.Printf("error received: %v\n", err)
			cancel()

		case <-sigIntChannel:
			if !aSigIntWasReceivedAlready {
				fmt.Println(ivkr.cliColorer.Error("Gracefully stopping... (signal Control-C again to force)"))
				aSigIntWasReceivedAlready = true
				cancel()
			} else {
				ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: "Terminated by Control-C", Code: 130})
				return // support fake exiter
			}

		case <-sigTermChannel:
			fmt.Println(ivkr.cliColorer.Error("Gracefully stopping..."))
			cancel()

		case event, isEventChannelOpen := <-eventChannel:
			if !isEventChannelOpen {
				ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1})
				return // support fake exiter
			}

			ivkr.cliOutput.Event(&event)

			if nil != event.CallEnded {
				if event.CallEnded.Call.ID == rootCallID {
					switch event.CallEnded.Outcome {
					case model.OpOutcomeSucceeded:
						ivkr.cliExiter.Exit(cliexiter.ExitReq{Code: 0})
					case model.OpOutcomeKilled:
						ivkr.cliExiter.Exit(cliexiter.ExitReq{Code: 137})
					default:
						// treat model.OpOutcomeFailed & unexpected values as errors.
						ivkr.cliExiter.Exit(cliexiter.ExitReq{Code: 1})
					}
					return // support fake exiter
				}
			}
		}
	}
}
