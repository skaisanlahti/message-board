package query

import "net/http"

func HomePage() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

	})
}
