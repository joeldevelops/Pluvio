package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joeldevelops/Pluvio/mdb"
)

type Config struct {
	Port string
	UsePhoneAuth bool
}

type API struct {
	app *fiber.App
	mongo *mdb.MongoDB
	config Config
}

// "/" endpoint
func (a *API) Index(c *fiber.Ctx) error {
	return c.SendString("You! Shall not! Pass!")
}

// "/health" endpoint
// Shuts down server if mongo cannot connect
func (a *API) HealthCheck(c *fiber.Ctx) error {
	log.Println("Health Check")
	ctx, ctxErr := context.WithTimeout(c.Context(), 30*time.Second)
	defer ctxErr()

	if ctxErr != nil {
		c.SendString("Context error in Heath Check")
	}

	if err := a.mongo.Ping(ctx, nil); err != nil {
		c.SendString("MongoDB is not connected")
		panic(1)
	}

	return c.SendString("OK")
}

// route initialization, additional routes in handler.go
func (a *API) setupRoutes() {
	a.app.Use(cors.New())

	a.app.Get("/", a.Index)
	a.app.Get("/health", a.HealthCheck)

	api := a.app.Group("/api")
	v1 := api.Group("/v1")
	// handler.go
	v1.Get("/rain/:timeRange", a.GetRainfallAmount)
	v1.Post("/rain", a.ReportRain)

	v1.Post("/user", a.CreateUser)
}

// Create an API instance and setup routes
func NewAPI(app *fiber.App, mongo *mdb.MongoDB, config Config) *API {
	a := &API{
		app,
		mongo,
		config,
	}

	a.setupRoutes()

	return a
}

// Start the server
func (a *API) StartServer() {
	log.Println("Starting server on port: " + a.config.Port)

	a.app.Listen(fmt.Sprintf(":%s", a.config.Port))
}
