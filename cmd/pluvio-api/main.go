package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joeldevelops/Pluvio/api"
	"github.com/joeldevelops/Pluvio/mdb"
	"github.com/joho/godotenv"
)

func getBoolEnv(key string) bool {
	return os.Getenv(key) == "true"
}

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}

	// Create new Fiber instance
	app := fiber.New()

	// Connect to MongoDB
	mongoConfig := &mdb.MDBConfig{}
	mongoConfig.DbName = os.Getenv("DB_NAME")
	mongoConfig.RainCollection = os.Getenv("DB_COLLECTION")
	mongoConfig.UsersCollection = os.Getenv("USERS_COLLECTION")

	log.Println("Connecting to MongoDB")
	mongo, err := mdb.NewMongoConnection(os.Getenv("MONGO_URL"), mongoConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Disconnect(context.TODO())

	// Create new config instance
	config := api.Config{
		Port: os.Getenv("PORT"),
		UsePhoneAuth: getBoolEnv("USE_PHONE_AUTH"),
	}

	// Create unique index for User collection concurrently. Idempoent.
	go func() {
		err := mongo.CreateUserIndex()
		if err != nil {
			log.Println("Error creating index for User collection")
			log.Println(err)
		}
	}()

	// Start server
	api.StartServer(app, mongo, config)
}