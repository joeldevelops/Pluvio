package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joeldevelops/Pluvio/mdb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertToVXML(message string) string {
	vxmlString := `
	<?xml version="1.0" ?>
	<vxml version="2.1">
		<form>
			<block>
				<prompt>%s</prompt>
			</block>
		</form>
	</vxml>
	`
	return fmt.Sprintf(vxmlString, message)
}

// CreateUser creates a new user in the database
func (a *API) CreateUser(c *fiber.Ctx) error {
	data := new(mdb.User)
	if err := c.BodyParser(data); err != nil {
		return err
	}

	log.Println("Creating user")
	oid, err := a.mongo.CreateUser(c.Context(), *data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.Status(409)
			return c.SendString("User already exists")
		}

		c.Status(500)
		return c.SendString("Error creating user")
	}

	c.Status(201)
	return c.SendString(fmt.Sprintf("Created user with ID: %s", oid.Hex()))
}

// GetRainfall returns the amount of rain reported in either the past day, week, or month
func (a *API) GetRainfallAmount(c *fiber.Ctx) error {
	timeRange := c.Params("timeRange", "day")
	location := c.Query("location", "")
	amount, err := a.calculateRainfall(c.Context(), timeRange, location)
	if err != nil {
		c.SendString(convertToVXML(fmt.Sprintf("Error getting %sly rainfall.", timeRange)))
	}

	message := fmt.Sprintf("In the past %s it rained %d milliliters.", timeRange, amount)

	c.Set("Content-Type", "application/xml")
	return c.SendString(convertToVXML(message))
}

// ReportRain allows a user to report rainfall in mm
func (a *API) ReportRain(c *fiber.Ctx) error {
	data := new(mdb.Rainfall)
	if err := c.BodyParser(data); err != nil {
		return err
	}

	if a.config.UsePhoneAuth {
		// Check if there is a phone number, if not return an error
		if data.PhoneNumber == "" {
			log.Println("no phone number provided")
			c.Status(400)
			return c.SendString(convertToVXML("No phone number provided, please call back and try again."))
		}

		// Check if the user exists in the DB
		if !a.mongo.CheckUserExists(c.Context(), data.PhoneNumber) {
			log.Println("Attempted use by unauthorized user")
			c.Status(403)
			return c.SendString(convertToVXML("You are not authorized to use this service."))
		}
		
		// Check if the user has already reported today
		if a.mongo.CheckUserReportedToday(c.Context(), data.PhoneNumber) {
			c.Status(429)
			return c.SendString(convertToVXML("Sorry, You have reached the maximum number of reports for today."))
		}
	}
		
	// Set the reportedAt field to the current time
	data.ReportedAt = primitive.NewDateTimeFromTime(time.Now())

	// INSERT
	oid, err := a.mongo.CreateRainfall(c.Context(), *data)
	if err != nil {
		return err
	}

	data.ID = oid
	log.Printf("Received rainfall report: %+v\n", data)

	// Return the ID of the inserted document
	// TODO: This should return success or a thank you message
	c.Status(201)
	return c.SendString(convertToVXML("Thank you for your report!"))
}

// Calculates the amount of rain based on the time range and location
func (a *API) calculateRainfall(ctx context.Context, timeRange string, loc string) (int, error) {
	// Set the filter based on the time range
	var tFilter time.Time
	switch timeRange {
	case "day":
		tFilter = time.Now().AddDate(0, 0, -1)
	case "week":
		tFilter = time.Now().AddDate(0, 0, -7)
	case "month":
		tFilter = time.Now().AddDate(0, -1, 0)
	default:
		return 0, fmt.Errorf("invalid time range")
	}

	// Set the filter to only return documents with a reportedAt field greater than the time filter
	// Could be day, week, month
	filter := bson.M{
		"reportedAt": bson.M{"$gte": primitive.NewDateTimeFromTime(tFilter)},
		"location":   bson.M{"$regex": primitive.Regex{Pattern: loc, Options: "i"}},
	}

	if loc == "" {
		delete(filter, "location")
	}

	results, err := a.mongo.GetRainfall(ctx, filter)
	if err != nil {
		return 0, err
	}

	// Iterate over the results and sum the amount of rain for the day
	var rain int
	for _, result := range results {
		rain += result.Amount
	}

	return rain, nil
}
