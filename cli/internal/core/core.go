package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"os"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// Core exposes all cli commands
//counterfeiter:generate -o fakes/core.go . Core
type Core interface {
	Auther
	Eventser
	Lser
	Oper
	Runer
	SelfUpdater
}

// New returns initialized cli core
func New(cliColorer clicolorer.CliColorer, containerRuntime, datadirPath string) Core {
	_os := ios.New()
	cliOutput := clioutput.New(cliColorer, os.Stderr, os.Stdout)
	cliExiter := cliexiter.New(cliOutput, _os)
	cliParamSatisfier := cliparamsatisfier.New(cliExiter, cliOutput)

	dataDir, err := datadir.New(datadirPath)
	if err != nil {
		panic(err)
	}

	if err := dataDir.InitAndLock(); nil != err {
		panic(err)
	}

	api := client.New(&client.Opts{
		ContainerRuntime: containerRuntime,
		DataDirPath:      dataDir.Path(),
	})

	dataResolver := dataresolver.New(
		cliExiter,
		cliParamSatisfier,
		api,
	)

	return _core{
		Auther: newAuther(
			cliExiter,
			dataResolver,
			api,
		),
		Eventser: newEventser(
			cliExiter,
			cliOutput,
			api,
		),
		Lser: newLser(
			cliExiter,
			cliOutput,
			dataResolver,
		),
		Oper: newOper(
			cliExiter,
			dataResolver,
			api,
		),
		Runer: newRuner(
			cliColorer,
			cliExiter,
			cliOutput,
			cliParamSatisfier,
			dataResolver,
			api,
		),
		SelfUpdater: newSelfUpdater(
			cliExiter,
			api,
		),
	}
}

type _core struct {
	Auther
	Eventser
	Lser
	Oper
	Runer
	SelfUpdater
}
