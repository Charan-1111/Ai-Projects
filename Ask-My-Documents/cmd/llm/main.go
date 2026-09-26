package main

import (
	"ask-my-documents/internal/server"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: could not load .env file:", err)
	}

	application, err := server.NewApplication()
	if err != nil {
		fmt.Println("Error initiating the server : ", err)
	}

	application.StartServer()
}
