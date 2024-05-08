package postgres

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/skaisanlahti/message-board/internal/config"
)

const (
	MigrateUp   = "up"
	MigrateDown = "down"
)

func Migrate(configuration config.Configuration, stdout io.Writer, database *sql.DB, direction string) error {
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
