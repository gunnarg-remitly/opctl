package main

import (
	"context"
	"fmt"
	"os"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	corePkg "github.com/opctl/opctl/cli/internal/core"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := newCli(ctx, corePkg.New)
	if err != nil {
		clicolorer.New().Error(fmt.Sprintf("failed to start up: %s", err.Error()))
		os.Exit(1)
	}
	cli.Run(os.Args)
}
