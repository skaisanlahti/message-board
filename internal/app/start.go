package app

import (
	"log"
	"net/http"

	"github.com/skaisanlahti/message-board/internal/home"
)

func Start() {
	router := http.NewServeMux()
	home.SetupRoutes(router)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("web server started on port 8080")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
