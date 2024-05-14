package command

import "net/http"

func SignUp() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

	})
}
