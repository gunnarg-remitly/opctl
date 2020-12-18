package main

import (
	"context"
	"os"

	corePkg "github.com/opctl/opctl/cli/internal/core"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newCli(ctx, corePkg.New).Run(os.Args)
}
