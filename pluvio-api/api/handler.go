package api

import (
	"context"
	"log"

	"github.com/joeldevelops/Pluvio/pluvio-api/mdb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) GetDayRain(c *fiber.Ctx) error {
	return c.SendString("10mm rain for the past day")
}

func (a *API) GetWeekRain(c *fiber.Ctx) error {
	return c.SendString("100mm rain for the past week")
}

func (a *API) GetMonthRain(c *fiber.Ctx) error {
	return c.SendString("1000mm rain for the past month")
}

func (a *API) ReportRain(c *fiber.Ctx) error {
	return c.SendString("Reported 10mm rain")
}

func (a *API) dbListRainfall(ctx context.Context) ([]*mdb.Rainfall, error) {
	findOptions := options.Find()
	
	collection := a.mongo.Database(a.config.DbName).Collection(a.config.DbCollection)

	cursor, err := collection.Find(ctx, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	var results []*mdb.Rainfall

	for cursor.Next(ctx) {
		var element mdb.Rainfall
		if err := cursor.Decode(&element); err != nil {
			log.Fatal(err)
			return nil, err
		}

		results = append(results, &element)
	}

	cursor.Close(ctx)

	return results, nil
}

func (a *API) dbWriteRainfall(data mdb.Rainfall) (primitive.ObjectID, error) {
	
	collection := a.mongo.Database(a.config.DbName).Collection(a.config.DbCollection)

	insertResult, err := collection.InsertOne(context.Background(), data)
	
	if err != nil {
		log.Fatal(err)
		return primitive.NilObjectID, err
	}

	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	} else {
		log.Fatal(err)
		return primitive.NilObjectID, err
	}
}