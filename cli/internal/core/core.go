package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/k8s"
)

// Core exposes all cli commands
//counterfeiter:generate -o fakes/core.go . Core
type Core interface {
	Auther
	Lser
	Oper
	Runer
	SelfUpdater
}

// New returns initialized cli core
func New(
	ctx context.Context,
	cliOutput clioutput.CliOutput,
	containerRuntime string,
	datadirPath string,
) Core {
	cliParamSatisfier := cliparamsatisfier.New(cliOutput)

	var cr containerruntime.ContainerRuntime
	var err error
	if "k8s" == containerRuntime {
		cr, err = k8s.New()
		if nil != err {
			panic(err)
		}
	} else {
		cr, err = docker.New(ctx)
		if nil != err {
			panic(err)
		}
	}

	eventChannel := make(chan model.Event)

	c, err := core.New(ctx, cr, datadirPath, eventChannel)
	if err != nil {
		panic(err)
	}

	dataResolver := dataresolver.New(
		cliParamSatisfier,
		c,
	)

	return _core{
		Auther: newAuther(
			dataResolver,
			c,
		),
		Lser: newLser(
			cliOutput,
			dataResolver,
		),
		Oper: newOper(
			dataResolver,
		),
		Runer: newRuner(
			cliOutput,
			cliParamSatisfier,
			dataResolver,
			eventChannel,
			c,
		),
		SelfUpdater: newSelfUpdater(),
	}
}

type _core struct {
	Auther
	Lser
	Oper
	Runer
	SelfUpdater
}
