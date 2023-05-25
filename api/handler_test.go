package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joeldevelops/Pluvio/mdb"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type TestAPI struct {
	*API
}

func (t *TestAPI) setup(usePhoneAuth bool) {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}

	testMdbConfig := &mdb.MDBConfig{
		DbName: fmt.Sprintf("%s_test", os.Getenv("DB_NAME")),
		RainCollection: fmt.Sprintf("%s_test", os.Getenv("DB_COLLECTION")),
		UsersCollection: fmt.Sprintf("%s_test", os.Getenv("USERS_COLLECTION")),
	}
	mongo, _ := mdb.NewMongoConnection(os.Getenv("MONGO_URL"), testMdbConfig)

	// Create new Fiber instance
	app := fiber.New()

	// Create new config instance
	config := Config{
		Port: os.Getenv("PORT"),
		UsePhoneAuth: usePhoneAuth,
	}

	// Create new API instance
	t.API = NewAPI(app, mongo, config)
}

func TestCreateUser(t *testing.T) {
	// setup
	a := &TestAPI{}
	a.setup(false)

	t.Run("CreateUser", func(t *testing.T) {
		// setup
		body := []byte(`{"name": "Joel", "phoneNumber": "+31612345678"}`)
		bodyReader := bytes.NewReader(body)
		req, _ := http.NewRequest("POST", "/api/v1/user", bodyReader)

		// test
		resp, err := a.API.app.Test(req)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})
}