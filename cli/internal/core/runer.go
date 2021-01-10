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
	"time"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	cliModel "github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/cli/internal/opstate"
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
		displayLiveGraph bool,
	) error
}

// newRuner returns an initialized "run" command
func newRuner(
	cliOutput clioutput.CliOutput,
	opFormatter clioutput.OpFormatter,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	dataResolver dataresolver.DataResolver,
	eventChannel chan model.Event,
	core core.Core,
) Runer {
	return _runer{
		cliOutput:         cliOutput,
		opFormatter:       opFormatter,
		cliParamSatisfier: cliParamSatisfier,
		dataResolver:      dataResolver,
		eventChannel:      eventChannel,
		core:              core,
	}
}

type _runer struct {
	dataResolver      dataresolver.DataResolver
	cliOutput         clioutput.CliOutput
	opFormatter       clioutput.OpFormatter
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier
	eventChannel      chan model.Event
	core              core.Core
}

func (ivkr _runer) Run(
	ctx context.Context,
	opRef string,
	opts *cliModel.RunOpts,
	displayLiveGraph bool,
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

	cliPromptInputSrc := ivkr.cliParamSatisfier.NewCliPromptInputSrc(opFile.Inputs)
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

	// listen for SIGINT on a channel
	aSigIntWasReceivedAlready := false
	sigIntChannel := make(chan os.Signal, 1)
	defer close(sigIntChannel)
	signal.Notify(
		sigIntChannel,
		syscall.SIGINT,
	)

	// listen for SIGTERM on a channel
	sigTermChannel := make(chan os.Signal, 1)
	defer close(sigTermChannel)
	signal.Notify(
		sigTermChannel,
		syscall.SIGTERM,
	)

	// listen for op end on a channel
	done := make(chan error, 1)
	go func() {
		_, err := ivkr.core.StartOp(
			ctx,
			model.StartOpReq{
				Args: argsMap,
				Op: model.StartOpReqOp{
					Ref: opHandle.Ref(),
				},
			},
		)
		done <- err
	}()

	// "request animation frame" like loop to force refresh of display loading spinners
	animationFrame := make(chan bool)
	if displayLiveGraph {
		go func() {
			for {
				time.Sleep(time.Second / 10)
				animationFrame <- true
			}
		}()
	}

	var state opstate.CallGraph
	var loadingSpinner opstate.DotLoadingSpinner
	output := opstate.NewOutputManager()

	defer func() {
		output.Print(state.String(ivkr.opFormatter, loadingSpinner, time.Now(), false))
		fmt.Println()
	}()

	clearGraph := func() {
		if displayLiveGraph {
			output.Clear()
		}
	}

	displayGraph := func() {
		if displayLiveGraph {
			output.Print(state.String(ivkr.opFormatter, loadingSpinner, time.Now(), true))
		}
	}

	for {
		select {
		case <-sigIntChannel:
			clearGraph()
			if !aSigIntWasReceivedAlready {
				ivkr.cliOutput.Warning("Gracefully stopping... (signal Control-C again to force)")
				aSigIntWasReceivedAlready = true
				// events will continue to stream in, make sure we continue to display the graph while this happens
				displayGraph()
				cancel()
			} else {
				return &RunError{
					ExitCode: 130,
					message:  "Terminated by Control-C",
				}
			}

		case <-sigTermChannel:
			clearGraph()
			ivkr.cliOutput.Error("Gracefully stopping...")
			displayGraph()
			cancel()

		case err := <-done:
			clearGraph()
			if !errors.Is(err, context.Canceled) {
				return err
			}
			return nil

		case event, isEventChannelOpen := <-ivkr.eventChannel:
			clearGraph()
			if !isEventChannelOpen {
				return errors.New("Event channel closed unexpectedly")
			}

			err := state.HandleEvent(&event)
			if err != nil {
				ivkr.cliOutput.Error(fmt.Sprintf("%v", err))
			}

			ivkr.cliOutput.Event(&event)
			displayGraph()
		case <-animationFrame:
			clearGraph()
			displayGraph()
		}
	}
}
