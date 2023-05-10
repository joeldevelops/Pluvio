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
	}

	// Create new config instance
	config := api.Config{
		DbName: os.Getenv("DB_NAME"),
		DbCollection: os.Getenv("DB_COLLECTION"),
		UserCollection: os.Getenv("USERS_COLLECTION"),
		Port: os.Getenv("PORT"),
	}

	// Create unique index for User collection concurrently. Idempoent.
	go func() {
		err := mdb.CreateUserIndex(mongo, config.DbName, config.UserCollection)
		if err != nil {
			log.Println("Error creating index for User collection")
			log.Println(err)
		}
	}()

	// Start server
	api.StartServer(app, mongo, config)
}