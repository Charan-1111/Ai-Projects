package main

import (
	"log"
	"tool-calling/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: could not load .env file:", err)
	}
	app, err := server.NewApplication()
	if err != nil {
		panic(err)
	}

	app.StartServer()
}
