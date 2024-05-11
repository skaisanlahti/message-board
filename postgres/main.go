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
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/message-board/library/file"
)

type appSettings struct {
	DatabaseAddress     string
	MigrationsDirectory string
}

const (
	MigrateUp   = "up"
	MigrateDown = "down"
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

	settingsFilePath := flag.String("settings", "", "application settings")
	direction := flag.String("migrate", "", "migration direction")
	flag.Parse()

	if *settingsFilePath == "" {
		return errors.New("settings not provided")
	}

	if *direction != MigrateUp && *direction != MigrateDown {
		return errors.New("invalid migration direction, must be either up or down")
	}

	config, err := file.ReadJSON[appSettings](*settingsFilePath)
	if err != nil {
		return err
	}

	database, err := sql.Open("pgx", config.DatabaseAddress)
	if err != nil {
		return err
	}
	defer database.Close()

	err = migrate(config, stdout, database, *direction)
	if err != nil {
		return err
	}

	return nil
}

func migrate(configuration appSettings, stdout io.Writer, database *sql.DB, direction string) error {
	directory := configuration.MigrationsDirectory
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}

	directory += direction
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		if direction == MigrateUp {
			return files[i].Name() < files[j].Name()
		}

		return files[i].Name() > files[j].Name()
	})

	fmt.Fprintf(stdout, "found %d migrations in %s\n", len(files), directory)
	for _, file := range files {
		migration, err := os.ReadFile(directory + "/" + file.Name())
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file.Name(), err)
		}
		_, err = database.Exec(string(migration))
		if err != nil {
			return fmt.Errorf("%s migration failed: %w", file.Name(), err)
		}
		fmt.Fprintf(stdout, "executed migration: %s\n", file.Name())
	}

	fmt.Fprintf(stdout, "migration %s finished\n", direction)
	return nil
}
