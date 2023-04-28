package main

import (
	"os"
	"log"

	"github.com/joeldevelops/Pluvio/pluvio-api/api"
	"github.com/joeldevelops/Pluvio/pluvio-api/mdb"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		log.Fatal(err)
	}

	app := fiber.New()

	log.Println("Connecting to MongoDB")
	mongo, err := mdb.ConnectMongo(os.Getenv("MONGO_URL"))

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	config := api.Config{
		DbName: os.Getenv("DB_NAME"),
		DbCollection: os.Getenv("DB_COLLECTION"),
		Port: os.Getenv("PORT"),
	}

	api.StartServer(app, mongo, config)
}