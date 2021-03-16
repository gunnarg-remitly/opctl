package core

import (
	"context"
	"fmt"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	callpkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/call"
)

//counterfeiter:generate -o internal/fakes/caller.go . caller
type caller interface {
	// Call executes a call
	Call(
		ctx context.Context,
		id string,
		scope map[string]*model.Value,
		callSpec *model.CallSpec,
		opPath string,
		parentCallID *string,
		rootCallID string,
	) (
		map[string]*model.Value,
		error,
	)
}

func newCaller(
	containerCaller containerCaller,
	dataDirPath string,
	eventChannel chan model.Event,
) caller {
	instance := &_caller{
		containerCaller: containerCaller,
		dataDirPath:     dataDirPath,
		eventChannel:    eventChannel,
	}
	instance.opCaller = newOpCaller(instance, dataDirPath)
	instance.parallelCaller = newParallelCaller(instance)
	instance.parallelLoopCaller = newParallelLoopCaller(instance)
	instance.serialCaller = newSerialCaller(instance)
	instance.serialLoopCaller = newSerialLoopCaller(instance)

	return instance
}

type _caller struct {
	containerCaller    containerCaller
	dataDirPath        string
	eventChannel       chan model.Event
	opCaller           opCaller
	parallelCaller     parallelCaller
	parallelLoopCaller parallelLoopCaller
	serialCaller       serialCaller
	serialLoopCaller   serialLoopCaller
}

func (clr _caller) Call(
	ctx context.Context,
	id string,
	scope map[string]*model.Value,
	callSpec *model.CallSpec,
	opPath string,
	parentCallID *string,
	rootCallID string,
) (
	map[string]*model.Value,
	error,
) {
	callCtx, cancelCall := context.WithCancel(ctx)
	defer cancelCall()
	var err error
	var outputs map[string]*model.Value
	var call *model.Call
	callStartTime := time.Now().UTC()

	if nil != callCtx.Err() {
		// if context done NOOP
		return nil, nil
	}

	// emit a call ended event after this call is complete
	defer func() {
		// defer must be defined before conditional return statements so it always runs

		if nil == call {
			call = &model.Call{
				ID:       id,
				RootID:   rootCallID,
				ParentID: parentCallID,
			}
		}

		event := model.Event{
			CallEnded: &model.CallEnded{
				Call:    *call,
				Outputs: outputs,
				Ref:     opPath,
			},
			Timestamp: time.Now().UTC(),
		}

		if nil != ctx.Err() {
			// this call or parent call killed/cancelled
			event.CallEnded.Outcome = model.OpOutcomeKilled
			event.CallEnded.Error = &model.CallEndedError{
				Message: ctx.Err().Error(),
			}
		} else if nil != err {
			event.CallEnded.Outcome = model.OpOutcomeFailed
			event.CallEnded.Error = &model.CallEndedError{
				Message: err.Error(),
			}
		} else {
			event.CallEnded.Outcome = model.OpOutcomeSucceeded
		}

		clr.eventChannel <- event
	}()

	if nil == callSpec {
		// NOOP
		return outputs, err
	}

	call, err = callpkg.Interpret(
		ctx,
		scope,
		callSpec,
		id,
		opPath,
		parentCallID,
		rootCallID,
		clr.dataDirPath,
	)
	if nil != err {
		return nil, err
	}

	if nil != call.If && !*call.If {
		return outputs, err
	}

	// Ensure this is emitted just after the deferred operation to emit the end
	// event is set up, so we always have a matching start and end event
	clr.eventChannel <- model.Event{
		Timestamp: callStartTime,
		CallStarted: &model.CallStarted{
			Call: *call,
			Ref:  opPath,
		},
	}

	switch {
	case nil != callSpec.Container:
		outputs, err = clr.containerCaller.Call(
			callCtx,
			call.Container,
			scope,
			callSpec.Container,
			rootCallID,
		)
	case nil != callSpec.Op:
		outputs, err = clr.opCaller.Call(
			callCtx,
			call.Op,
			scope,
			parentCallID,
			rootCallID,
			callSpec.Op,
		)
	case nil != callSpec.Parallel:
		outputs, err = clr.parallelCaller.Call(
			callCtx,
			id,
			scope,
			rootCallID,
			opPath,
			*callSpec.Parallel,
		)
	case nil != callSpec.ParallelLoop:
		outputs, err = clr.parallelLoopCaller.Call(
			callCtx,
			id,
			scope,
			*callSpec.ParallelLoop,
			opPath,
			parentCallID,
			rootCallID,
		)
	case nil != callSpec.Serial:
		outputs, err = clr.serialCaller.Call(
			callCtx,
			id,
			scope,
			rootCallID,
			opPath,
			*callSpec.Serial,
		)
	case nil != callSpec.SerialLoop:
		outputs, err = clr.serialLoopCaller.Call(
			callCtx,
			id,
			scope,
			*callSpec.SerialLoop,
			opPath,
			parentCallID,
			rootCallID,
		)
	default:
		err = fmt.Errorf("invalid call graph '%+v'", callSpec)
	}

	return outputs, err
}
