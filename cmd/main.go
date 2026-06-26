package main

import (
	"Task-Management/database"
	"Task-Management/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	if err := database.ConnectDB(); err != nil {
		log.Panicf("failed to connect to DB with err: %v", err)
	}

	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Panicf("failed to close DB connection with error: %v", err)
		}

		// close server
	}()

	r := server.SetupRoutes()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Panicf("failed to start server with err: %v", err)
	}

	fmt.Println("Listening on port 8080")
}
