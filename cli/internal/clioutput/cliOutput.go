package clioutput

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/sdks/go/model"
)

//CliOutput allows mocking/faking output
//counterfeiter:generate -o fakes/cliOutput.go . CliOutput
type CliOutput interface {
	// silently disables coloring
	DisableColor()

	// outputs a msg requiring attention
	Attention(s string)

	// outputs a warning message (looks like an error but on stdout)
	Warning(s string)

	// outputs an error msg
	Error(s string)

	// outputs an event
	// @TODO: not generic
	Event(event *model.Event)

	// outputs a success msg
	Success(s string)
}

func New(
	cliColorer clicolorer.CliColorer,
	datadirPath string,
	errWriter io.Writer,
	stdWriter io.Writer,
) (CliOutput, error) {
	return _cliOutput{
		cliColorer:  cliColorer,
		datadirPath: datadirPath,
		errWriter:   errWriter,
		stdWriter:   stdWriter,
	}, nil
}

type _cliOutput struct {
	cliColorer  clicolorer.CliColorer
	datadirPath string
	errWriter   io.Writer
	stdWriter   io.Writer
}

func (this _cliOutput) DisableColor() {
	this.cliColorer.DisableColor()
}

func (this _cliOutput) Attention(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Attention(s),
		),
	)
}

func (this _cliOutput) Warning(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Error(s),
		),
	)
}

func (this _cliOutput) Error(s string) {
	io.WriteString(
		this.errWriter,
		fmt.Sprintln(
			this.cliColorer.Error(s),
		),
	)
}

func (this _cliOutput) Event(event *model.Event) {
	switch {
	case nil != event.CallEnded &&
		nil == event.CallEnded.Call.Op &&
		nil == event.CallEnded.Call.Container &&
		nil != event.CallEnded.Error:
		this.error(event)

	case nil != event.CallEnded &&
		nil != event.CallEnded.Call.Container:
		this.containerExited(event)

	case nil != event.CallStarted &&
		nil != event.CallStarted.Call.Container:
		this.containerStarted(event)

	case nil != event.ContainerStdErrWrittenTo:
		this.containerStdErrWrittenTo(event.ContainerStdErrWrittenTo)

	case nil != event.ContainerStdOutWrittenTo:
		this.containerStdOutWrittenTo(event.ContainerStdOutWrittenTo)

	case nil != event.CallEnded &&
		nil != event.CallEnded.Call.Op:
		this.opEnded(event)

	case nil != event.CallStarted &&
		nil != event.CallStarted.Call.Op:
		this.opStarted(event.CallStarted)
	}
}

func (this _cliOutput) error(event *model.Event) {
	io.WriteString(
		this.errWriter,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.CallEnded.Call.ID, event.CallEnded.Ref),
			this.cliColorer.Error(event.CallEnded.Error.Message),
		),
	)
}

func (this _cliOutput) containerExited(event *model.Event) {
	var color func(s string) string
	var writer io.Writer
	switch event.CallEnded.Outcome {
	case model.OpOutcomeSucceeded:
		color = this.cliColorer.Success
		writer = this.stdWriter
	case model.OpOutcomeKilled:
		color = this.cliColorer.Info
		writer = this.stdWriter
	default:
		color = this.cliColorer.Error
		writer = this.errWriter
	}

	message := "exited"
	if nil != event.CallEnded.Error {
		message += fmt.Sprintf(" with error")
	}
	message = color(message)
	if nil != event.CallEnded.Call.Container.Image.Ref {
		message += fmt.Sprintf(" %s", *event.CallEnded.Call.Container.Image.Ref)
	}
	if nil != event.CallEnded.Error {
		message += fmt.Sprintf("\n%v", event.CallEnded.Error.Message)
	}

	io.WriteString(
		writer,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.CallEnded.Call.ID, event.CallEnded.Ref),
			message,
		),
	)
}

func (this _cliOutput) containerStarted(event *model.Event) {
	message := this.cliColorer.Info("started container")
	if nil != event.CallStarted.Call.Container.Image.Ref {
		message += fmt.Sprintf(" %s", *event.CallStarted.Call.Container.Image.Ref)
	}

	io.WriteString(
		this.stdWriter,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.CallStarted.Call.ID, event.CallStarted.Ref),
			message,
		),
	)
}

func (this _cliOutput) formatOpRef(opRef string) string {
	if path.IsAbs(opRef) {
		cwd, err := os.Getwd()
		if err != nil {
			return opRef
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return opRef
		}
		dataDirPath := this.datadirPath
		if strings.HasPrefix(opRef, dataDirPath) {
			return opRef[len(dataDirPath+string(os.PathListSeparator)+"ops"+string(os.PathListSeparator)):]
		}
		if strings.HasPrefix(opRef, cwd) {
			return "." + opRef[len(cwd):]
		}
		if strings.HasPrefix(opRef, home) {
			return "~" + opRef[len(home):]
		}
	}
	return opRef
}

func (this _cliOutput) outputPrefix(id, opRef string) string {
	parts := []string{id[:8]}
	opRef = this.formatOpRef(opRef)
	if opRef != "" {
		parts = append(parts, opRef)
	}
	return this.cliColorer.Muted(strings.Join(parts, " ")) + ": "
}

func (this _cliOutput) containerStdErrWrittenTo(event *model.ContainerStdErrWrittenTo) {
	io.WriteString(
		this.errWriter,
		fmt.Sprintf(
			"%s%s",
			this.outputPrefix(event.ContainerID, event.OpRef),
			event.Data,
		),
	)
}

func (this _cliOutput) containerStdOutWrittenTo(event *model.ContainerStdOutWrittenTo) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintf(
			"%s%s",
			this.outputPrefix(event.ContainerID, event.OpRef),
			event.Data,
		),
	)
}

func (this _cliOutput) opEnded(event *model.Event) {
	var color func(s string) string
	var writer io.Writer
	switch event.CallEnded.Outcome {
	case model.OpOutcomeSucceeded:
		color = this.cliColorer.Success
		writer = this.stdWriter
	case model.OpOutcomeKilled:
		color = this.cliColorer.Info
		writer = this.stdWriter
	default:
		color = this.cliColorer.Error
		writer = this.errWriter
	}

	message := "ended"
	if nil != event.CallEnded.Error {
		message += fmt.Sprintf(" with error")
	}
	message = color(message)
	if nil != event.CallEnded.Error {
		message += fmt.Sprintf("\n%v", event.CallEnded.Error.Message)
	}

	io.WriteString(
		writer,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.CallEnded.Call.ID, event.CallEnded.Call.Op.OpPath),
			message,
		),
	)
}

func (this _cliOutput) opStarted(event *model.CallStarted) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.Call.ID, event.Call.Op.OpPath),
			this.cliColorer.Info("started op"),
		),
	)
}

func (this _cliOutput) info(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Info(s),
		),
	)
}

func (this _cliOutput) Success(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Success(s),
		),
	)
}
