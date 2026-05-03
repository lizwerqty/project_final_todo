package main

import (
	"log"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	server.Run()
}
