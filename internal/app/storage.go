package app

import (
	"database/sql"

	"github.com/skaisanlahti/message-board/internal/assert"
	"github.com/skaisanlahti/message-board/internal/core"
)

type AppStorage struct {
	db *sql.DB
}

func (s *AppStorage) Database() *sql.DB {
	return s.db
}

func NewAppStorage(config core.Configuration) *AppStorage {
	db, err := sql.Open("pgx", config.DatabaseAddress())
	assert.Ok(err, "failed to open database")
	return &AppStorage{db}
}
