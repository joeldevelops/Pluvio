package api

import (
	"context"
	"log"
	"os"
	"time"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	DbName string
	DbCollection string
	Port string
}

type API struct {
	app *fiber.App
	mongo *mongo.Client
	config Config
}

func (a *API) Index(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (a *API) HealthCheck(c *fiber.Ctx) error {
	log.Println("Health Check")
	ctx, ctxErr := context.WithTimeout(c.Context(), 30*time.Second)
	defer ctxErr()

	if ctxErr != nil {
		c.SendString("Context error in Heath Check")
	}

	if err := a.mongo.Ping(ctx, nil); err != nil {
		c.SendString("MongoDB is not connected")
		os.Exit(1)
	}

	return c.SendString("OK")
}

func (a *API) setupRoutes() {
	a.app.Get("/", a.Index)
	a.app.Get("/health", a.HealthCheck)

	a.app.Get("api/v1/rain/day", a.GetDayRain)
	a.app.Get("api/v1/rain/week", a.GetWeekRain)
	a.app.Get("api/v1/rain/month", a.GetMonthRain)
	a.app.Post("api/v1/rain", a.ReportRain)
}

func StartServer(app *fiber.App, mongo *mongo.Client, config Config) {
	a := &API{
		app: app,
		mongo: mongo,
		config: config,
	}

	a.setupRoutes()

	log.Println("Starting server on port: " + a.config.Port)

	app.Listen(fmt.Sprintf(":%s", a.config.Port))
}
