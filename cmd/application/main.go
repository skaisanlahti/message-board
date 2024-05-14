package main

import (
	"context"
	"fmt"
	"os"

	"github.com/skaisanlahti/message-board/internal/application"
)

func main() {
	ctx := context.Background()
	err := application.Run(ctx, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "main error: %s\n", err)
		os.Exit(1)
	}
}
