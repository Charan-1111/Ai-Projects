package main

import (
	"llm-playground/internal/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// this is the entry point of the appliaction
	// responsible only for application startup

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: could not load .env file:", err)
	}
	application, err := server.NewApplication()
	if err != nil {
		log.Println("Error loading the server : " + err.Error())
		os.Exit(1)
	}

	err = application.StartServer()
	if err != nil {
		log.Println("Error starting the server : " + err.Error())
		os.Exit(1)
	}
}
