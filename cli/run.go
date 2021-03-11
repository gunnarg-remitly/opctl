package main

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
	"github.com/opctl/opctl/cli/internal/opgraph"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// RunOpts are options to run a given op through the CLI
type RunOpts struct {
	ArgFile string
	Args    []string
}

func run(
	ctx context.Context,
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	eventChannel chan model.Event,
	node core.Core,
	opFormatter clioutput.OpFormatter,
	opRef string,
	opts *RunOpts,
	displayLiveGraph bool,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dataResolver := dataresolver.New(
		cliParamSatisfier,
		node,
	)

	opHandle, err := dataResolver.Resolve(
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

	ymlFileInputSrc, err := cliParamSatisfier.NewYMLFileInputSrc(opts.ArgFile)
	if nil != err {
		return fmt.Errorf("unable to load arg file at '%v'; error was: %v", opts.ArgFile, err.Error())
	}

	cliPromptInputSrc := cliParamSatisfier.NewCliPromptInputSrc(opFile.Inputs)
	if nil != err {
		return err
	}
	argsMap, err := cliParamSatisfier.Satisfy(
		cliparamsatisfier.NewInputSourcer(
			cliParamSatisfier.NewSliceInputSrc(opts.Args, "="),
			ymlFileInputSrc,
			cliParamSatisfier.NewEnvVarInputSrc(),
			cliParamSatisfier.NewParamDefaultInputSrc(opFile.Inputs),
			cliPromptInputSrc,
		),
		opFile.Inputs,
	)
	if nil != err {
		return err
	}

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

	sigInfoChannel := make(chan os.Signal, 1)
	defer close(sigInfoChannel)
	signal.Notify(
		sigInfoChannel,
		syscall.Signal(0x1d), // portable version of syscall.SIGINFO
	)

	// listen for op end on a channel
	done := make(chan error, 1)
	go func() {
		_, err := node.StartOp(
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

	var state opgraph.CallGraph
	var loadingSpinner opgraph.DotLoadingSpinner
	output := opgraph.NewOutputManager()

	defer func() {
		output.Print(state.String(opFormatter, loadingSpinner, time.Now(), false))
		fmt.Println()
	}()

	clearGraph := func() {
		if displayLiveGraph {
			output.Clear()
		}
	}

	displayGraph := func() {
		if displayLiveGraph {
			output.Print(state.String(opFormatter, loadingSpinner, time.Now(), true))
		}
	}

	for {
		select {
		case <-sigIntChannel:
			clearGraph()
			if !aSigIntWasReceivedAlready {
				cliOutput.Warning("Gracefully stopping... (signal Control-C again to force)")
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

		case <-sigInfoChannel:
			clearGraph()
			// clear two more lines
			fmt.Print("\033[1A\033[K\033[1A\033[K")
			fmt.Println(state.String(opFormatter, opgraph.StaticLoadingSpinner{}, time.Now(), false))
			displayGraph()

		case <-sigTermChannel:
			clearGraph()
			cliOutput.Error("Gracefully stopping...")
			displayGraph()
			cancel()

		case err := <-done:
			clearGraph()
			if !errors.Is(err, context.Canceled) {
				return err
			}
			return nil

		case event, isEventChannelOpen := <-eventChannel:
			clearGraph()
			if !isEventChannelOpen {
				return errors.New("Event channel closed unexpectedly")
			}

			if err := state.HandleEvent(&event); err != nil {
				cliOutput.Error(fmt.Sprintf("%v", err))
			}

			cliOutput.Event(&event)
			displayGraph()
		case <-animationFrame:
			clearGraph()
			displayGraph()
		}
	}
}
