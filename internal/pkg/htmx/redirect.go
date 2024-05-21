package htmx

import "net/http"

func Redirect(response http.ResponseWriter, url string) {
	response.Header().Add("HX-Location", url)
	response.WriteHeader(http.StatusOK)
}
