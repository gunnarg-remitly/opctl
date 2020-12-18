package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	corePkg "github.com/opctl/opctl/cli/internal/core"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := newCli(ctx, corePkg.New)
	if err != nil {
		clicolorer.New().Error(fmt.Sprintf("failed to start up: %v", err.Error()))
		os.Exit(1)
	}

	defer func() {
		if panicArg := recover(); panicArg != nil {
			clicolorer.New().Error(fmt.Sprintf("panic: %v", panicArg))
			fmt.Println(panicArg, debug.Stack())
		}
	}()

	cli.Run(os.Args)
}
