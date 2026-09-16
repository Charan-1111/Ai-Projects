package main

import (
	"log"
	"os"

	"semantic-search/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: could not load .env file:", err)
	}

	application, err := server.NewApplication()
	if err != nil {
		log.Println("Error loading the server:", err)
		os.Exit(1)
	}

	if err := application.StartServer(); err != nil {
		log.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
