package core

import "net/http"

type Router interface {
	Route(pattern string, handler http.HandlerFunc)
}
