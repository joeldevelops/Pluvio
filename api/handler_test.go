package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/joeldevelops/Pluvio/mdb"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestAPI struct {
	*API
}

func (t *TestAPI) setup(usePhoneAuth bool) {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}

	testMdbConfig := &mdb.MDBConfig{
		DbName: fmt.Sprintf("%s_test", os.Getenv("DB_NAME")),
		RainCollection: fmt.Sprintf("%s_test", os.Getenv("DB_COLLECTION")),
		UsersCollection: fmt.Sprintf("%s_test", os.Getenv("USERS_COLLECTION")),
	}
	url := os.Getenv("MONGO_URL")
	db, _ := mdb.NewMongoConnection(url, testMdbConfig)

	// Create new Fiber instance
	app := fiber.New()

	// Create new config instance
	config := Config{
		Port: os.Getenv("PORT"),
		UsePhoneAuth: usePhoneAuth,
	}
	
	// Create unique index for test User collection concurrently. Idempoent.
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"phoneNumber": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	// Index creation
	_, _ = db.Database(testMdbConfig.DbName).Collection(testMdbConfig.UsersCollection).Indexes().CreateOne(context.TODO(), indexModel)

	// Create new API instance
	t.API = NewAPI(app, db, config)
}

func (t *TestAPI) teardown() {
	// Drop test database
	testDbName := fmt.Sprintf("%s_test", os.Getenv("DB_NAME"))
	t.API.mongo.Database(testDbName).Drop(context.Background())

	// Disconnect from MongoDB
	t.API.mongo.Disconnect(context.Background())
}

func TestCreateUser(t *testing.T) {
	// setup
	a := &TestAPI{}
	a.setup(false)
	defer a.teardown()

	t.Run("Should create a user", func(t *testing.T) {
		// setup
		body := mdb.User{Name: "joel", PhoneNumber: "1234567890"}
		bodyJSON, _ := json.Marshal(&body)
		req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		// test
		resp, err := a.API.app.Test(req, -1)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("Should error on duplicate user", func(t *testing.T) {
		// setup
		body := mdb.User{Name: "joel", PhoneNumber: "1234567890"}
		bodyJSON, _ := json.Marshal(&body)
		req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		// test
		resp, err := a.API.app.Test(req, -1)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, 409, resp.StatusCode)
	})
}