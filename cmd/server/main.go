package main

import (
	"context"
	"fmt"
	"os"

	"github.com/skaisanlahti/message-board/internal/app"
)

func main() {
	ctx := context.Background()
	err := app.Run(ctx, os.Args, os.Getenv, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
