package main

import (
	"os"
	"log"

	"github.com/joeldevelops/Pluvio/api"
	"github.com/joeldevelops/Pluvio/mdb"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}

	// Create new Fiber instance
	app := fiber.New()

	log.Println("Connecting to MongoDB")
	mongo, err := mdb.ConnectMongo(os.Getenv("MONGO_URL"))

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Create new config instance
	config := api.Config{
		DbName: os.Getenv("DB_NAME"),
		DbCollection: os.Getenv("DB_COLLECTION"),
		Port: os.Getenv("PORT"),
	}

	// Start server
	api.StartServer(app, mongo, config)
}