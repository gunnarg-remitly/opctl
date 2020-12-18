package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"os"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/dataresolver"
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
func New(ctx context.Context, cliColorer clicolorer.CliColorer, containerRuntime, datadirPath string) Core {
	dataDir, err := datadir.New(datadirPath)
	if err != nil {
		panic(err)
	}
	if err := dataDir.InitAndLock(); nil != err {
		panic(err)
	}

	_os := ios.New()
	cliOutput := clioutput.New(cliColorer, dataDir.Path(), os.Stderr, os.Stdout)
	cliExiter := cliexiter.New(cliOutput, _os)
	cliParamSatisfier := cliparamsatisfier.New(cliExiter, cliOutput)

	var cr containerruntime.ContainerRuntime
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
	c, err := core.New(ctx, cr, datadirPath)
	if err != nil {
		panic(err)
	}

	dataResolver := dataresolver.New(
		cliExiter,
		cliParamSatisfier,
		c,
	)

	return _core{
		Auther: newAuther(
			cliExiter,
			dataResolver,
			c,
		),
		Lser: newLser(
			cliExiter,
			cliOutput,
			dataResolver,
		),
		Oper: newOper(
			cliExiter,
			dataResolver,
		),
		Runer: newRuner(
			cliColorer,
			cliExiter,
			cliOutput,
			cliParamSatisfier,
			dataResolver,
			c,
		),
		SelfUpdater: newSelfUpdater(cliExiter),
	}
}

type _core struct {
	Auther
	Lser
	Oper
	Runer
	SelfUpdater
}
