package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/message-board/internal/postgres"
)

func main() {
	ctx := context.Background()
	err := postgres.Run(ctx, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
