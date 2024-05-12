package main

import (
	"context"
	"fmt"
	"os"

	"github.com/skaisanlahti/message-board/application/program"
)

func main() {
	ctx := context.Background()
	err := program.Run(ctx, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "main error: %s\n", err)
		os.Exit(1)
	}
}
