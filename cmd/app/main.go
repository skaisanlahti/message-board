package main

import (
	"context"
	"fmt"
	"os"

	"github.com/skaisanlahti/message-board/internal/app"
)

func main() {
	ctx := context.Background()
	err := app.Run(ctx, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "main error: %s\n", err)
		os.Exit(1)
	}
}
