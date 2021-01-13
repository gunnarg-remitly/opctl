package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/lazylocalnode"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/runtime"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/core"
)

// Core exposes all cli commands
//counterfeiter:generate -o fakes/core.go . Core
type Core interface {
	Auther
	Eventser
	Lser
	Noder
	Oper
	Runer
	SelfUpdater
	UIer
}

// New returns initialized cli core
func New(
	cliOutput clioutput.CliOutput,
	nodeProviderOpts nodeprovider.NodeOpts,
) Core {
	cliParamSatisfier := cliparamsatisfier.New(cliOutput)

	var nodeProvider nodeprovider.NodeProvider
	var opNode node.OpNode

	if nodeProviderOpts.DisableNode {
		containerRuntime, _ := runtime.New(nodeProviderOpts.ContainerRuntime)
		opNode = core.New(
			containerRuntime,
			nodeProviderOpts.DataDir,
		)
	} else {
		nodeProvider = local.New(nodeProviderOpts)
		opNode = lazylocalnode.New(nodeProvider)
	}

	dataResolver := dataresolver.New(
		cliParamSatisfier,
		opNode,
	)

	return _core{
		Auther: newAuther(
			dataResolver,
			opNode,
		),
		Eventser: newEventser(
			cliOutput,
			opNode,
		),
		Lser: newLser(
			cliOutput,
			dataResolver,
		),
		Noder: newNoder(nodeProvider),
		Oper: newOper(
			dataResolver,
			opNode,
		),
		Runer: newRuner(
			cliOutput,
			cliParamSatisfier,
			dataResolver,
			opNode,
		),
		SelfUpdater: newSelfUpdater(nodeProvider),
		UIer: newUIer(
			dataResolver,
			opNode,
		),
	}
}

type _core struct {
	Auther
	Eventser
	Lser
	Noder
	Oper
	Runer
	SelfUpdater
	UIer
}
