package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Aman913k/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Println("Warning: MONGO_URI not set, using default MongoDB URI")
		mongoURI = "mongodb://localhost:27017" // Set your default URI here if needed
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Error: JWT_SECRET environment variable is required")
	}

	// Continue with application setup...
	log.Println("Mongo URI:", mongoURI)
	log.Println("JWT Secret loaded successfully.")

	r := routes.Router()
	fmt.Println("Server is getting Started...")

	log.Fatal(http.ListenAndServe(":5000", r))
	fmt.Println("Listening at port 5000...")
}
