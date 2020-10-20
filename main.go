package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Loading .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	a := App{}
	a.Initialize(os.Getenv("DB_NAME"))
	a.Run("127.0.0.1:8080")
}
