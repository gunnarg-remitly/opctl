package core

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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
	) error
}

// newRuner returns an initialized "run" command
func newRuner(
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	datadirPath string,
	dataResolver dataresolver.DataResolver,
	eventChannel chan model.Event,
	core core.Core,
) Runer {
	return _runer{
		cliOutput:         cliOutput,
		cliParamSatisfier: cliParamSatisfier,
		datadirPath:       datadirPath,
		dataResolver:      dataResolver,
		eventChannel:      eventChannel,
		core:              core,
	}
}

type _runer struct {
	dataResolver      dataresolver.DataResolver
	datadirPath       string
	cliOutput         clioutput.CliOutput
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier
	eventChannel      chan model.Event
	core              core.Core
}

func (ivkr _runer) Run(
	ctx context.Context,
	opRef string,
	opts *cliModel.RunOpts,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opHandle, err := ivkr.dataResolver.Resolve(
		ctx,
		opRef,
		nil,
	)
	if nil != err {
		return err
	}

	opFileReader, err := opHandle.GetContent(
		ctx,
		opfile.FileName,
	)
	if nil != err {
		return err
	}

	opFileBytes, err := ioutil.ReadAll(opFileReader)
	if nil != err {
		return err
	}

	opFile, err := opfile.Unmarshal(
		filepath.Join(opHandle.Ref(), opfile.FileName),
		opFileBytes,
	)
	if nil != err {
		return err
	}

	ymlFileInputSrc, err := ivkr.cliParamSatisfier.NewYMLFileInputSrc(opts.ArgFile)
	if nil != err {
		return fmt.Errorf("unable to load arg file at '%v'; error was: %v", opts.ArgFile, err.Error())
	}

	cliPromptInputSrc, err := ivkr.cliParamSatisfier.NewCliPromptInputSrc(ivkr.datadirPath, opFile.Inputs)
	if nil != err {
		return err
	}
	argsMap, err := ivkr.cliParamSatisfier.Satisfy(
		cliparamsatisfier.NewInputSourcer(
			ivkr.cliParamSatisfier.NewSliceInputSrc(opts.Args, "="),
			ymlFileInputSrc,
			ivkr.cliParamSatisfier.NewEnvVarInputSrc(),
			ivkr.cliParamSatisfier.NewParamDefaultInputSrc(opFile.Inputs),
			cliPromptInputSrc,
		),
		opFile.Inputs,
	)
	if nil != err {
		return err
	}

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
		return err
	}

	for {
		select {
		case <-sigIntChannel:
			if !aSigIntWasReceivedAlready {
				ivkr.cliOutput.Warning("Gracefully stopping... (signal Control-C again to force)")
				aSigIntWasReceivedAlready = true
				cancel()
			} else {
				return &RunError{
					ExitCode: 130,
					message:  "Terminated by Control-C",
				}
			}

		case <-sigTermChannel:
			ivkr.cliOutput.Error("Gracefully stopping...")
			cancel()

		case event, isEventChannelOpen := <-ivkr.eventChannel:
			if !isEventChannelOpen {
				return errors.New("Event channel closed unexpectedly")
			}

			if ctx.Err() != context.Canceled {
				ivkr.cliOutput.Event(&event)
			}

			if nil != event.CallEnded {
				if event.CallEnded.Call.ID == rootCallID {
					switch event.CallEnded.Outcome {
					case model.OpOutcomeSucceeded:
						return nil
					case model.OpOutcomeKilled:
						return &RunError{ExitCode: 137}
					default:
						return &RunError{ExitCode: 1}
					}
				}
			}
		}
	}
}
