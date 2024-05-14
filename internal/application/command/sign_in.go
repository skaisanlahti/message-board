package command

import "net/http"

func SignIn() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

	})
}
