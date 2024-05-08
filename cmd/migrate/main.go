package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/message-board/internal/config"
	"github.com/skaisanlahti/message-board/internal/postgres"
)

func main() {
	ctx := context.Background()
	err := run(ctx, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, stdout io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	direction := flag.String("m", "", "")
	flag.Parse()
	if *direction != postgres.MigrateUp && *direction != postgres.MigrateDown {
		return errors.New("invalid migration direction, must be either up or down")
	}

	configuration, err := config.Read("config/app.json")
	if err != nil {
		return err
	}

	database, err := sql.Open("pgx", configuration.DatabaseAddress)
	if err != nil {
		return err
	}
	defer database.Close()

	err = postgres.Migrate(configuration, stdout, database, *direction)
	if err != nil {
		return err
	}

	return nil
}
