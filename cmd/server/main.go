package main

import (
	"context"
	"fmt"
	"os"

	"github.com/skaisanlahti/message-board/internal/app"
)

func main() {
	runCtx := context.Background()
	err := app.Run(runCtx, os.Args, os.Getenv, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "app.Run error: %s\n", err)
		os.Exit(1)
	}
}
