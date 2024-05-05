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

	direction := flag.String("cmd", "", "")
	flag.Parse()
	if *direction == "" {
		return errors.New("direction argument was not provided")
	}

	configuration, err := config.Read("appsettings.json")
	if err != nil {
		return err
	}

	database, err := sql.Open("pgx", configuration.DatabaseAddress)
	if err != nil {
		return err
	}
	defer database.Close()

	err = postgres.RunMigrations(configuration, stdout, database, *direction)
	if err != nil {
		return err
	}

	return nil
}