package web

import "net/http"

func ServeStaticFiles() {
	handler := http.FileServerFS()
}

func ParseTemplates() {

}
