package migrator

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"io"
	"os"
	"os/signal"

	"github.com/skaisanlahti/message-board/internal/file"
)

type appSettings struct {
	DatabaseAddress     string
	MigrationsDirectory string
}

const (
	MigrateUp   = "up"
	MigrateDown = "down"
)

func Run(ctx context.Context, stdout io.Writer) error {
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

	err = migrate(config.MigrationsDirectory, stdout, database, *direction)
	if err != nil {
		return err
	}

	return nil
}
