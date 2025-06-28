package main

import (
	"GoChain/server"
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
