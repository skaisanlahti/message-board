package program

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func migrate(migrationsDirectory string, stdout io.Writer, database *sql.DB, direction string) error {
	directory := migrationsDirectory
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	matcher := direction + ".sql"
	migrations := []os.DirEntry{}
	for _, file := range files {
		if strings.Contains(file.Name(), matcher) {
			migrations = append(migrations, file)
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		if direction == MigrateUp {
			return migrations[i].Name() < migrations[j].Name()
		}

		return migrations[i].Name() > migrations[j].Name()
	})

	fmt.Fprintf(stdout, "found %d matching migrations in %s\n", len(migrations), directory)
	for _, migration := range migrations {
		statement, err := os.ReadFile(directory + "/" + migration.Name())
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", migration.Name(), err)
		}
		_, err = database.Exec(string(statement))
		if err != nil {
			return fmt.Errorf("%s migration failed: %w", migration.Name(), err)
		}
		fmt.Fprintf(stdout, "executed migration: %s\n", migration.Name())
	}

	fmt.Fprintf(stdout, "migration %s finished\n", direction)
	return nil
}
