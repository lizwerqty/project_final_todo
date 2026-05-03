package server

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

func Run() {
	port := "7540"

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	api.Init()

	webDir := "./web"

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Server started on :%s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
