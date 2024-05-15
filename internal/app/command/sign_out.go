package command

import "net/http"

func SignOut() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

	})
}
