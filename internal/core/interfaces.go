package core

import (
	"database/sql"
	"net/http"
)

type Configuration interface {
	ServerAddress() string
	DatabaseAddress() string
}

type Server interface {
	Start()
}

type Router interface {
	Handle(pattern string, handler http.HandlerFunc)
	Handler() http.Handler
}

type Storage interface {
	Database() *sql.DB
}

type Renderer interface {
	Render(name string, data any) []byte
}
