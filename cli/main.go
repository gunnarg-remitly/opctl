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
	cliColorer := clicolorer.New()
	defer func() {
		if panicArg := recover(); panicArg != nil {
			cancel()
			fmt.Println(
				cliColorer.Error(
					fmt.Sprint(panicArg),
				),
			)
			os.Exit(1)
		}
	}()

	newCli(
		ctx,
		cliColorer,
		corePkg.New,
	).
		Run(os.Args)

}
