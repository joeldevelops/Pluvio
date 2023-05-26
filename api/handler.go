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

func convertToVXML(message string, language string) string {
	if language == "french" {
		language = "fr-FR"
	} else {
		language = "en-US"
	}

	vxmlString := `
	<?xml version="1.0" ?>
	<vxml version="2.1" xml:lang="%s">
		<form>
			<block>
				<prompt>"%s"</prompt>
			</block>
		</form>
	</vxml>
	`
	return fmt.Sprintf(vxmlString, language, message)
}

func buildResponse(ctx *fiber.Ctx, message string, rainfall int, lang string) string {
	if ctx.Accepts("application/xml") != "" {
		return convertToVXML(message, lang)
	}

	if ctx.Accepts("application/json", "json") != "" {
		if rainfall < 0 {
			return fmt.Sprintf(`{"message": "%s"}`, message)
		}
		return fmt.Sprintf(`{"message": "%s", "rainfall":%d}`, message, rainfall)
	}

	// default to xml
	return convertToVXML(message, lang)
}

// Send the response as XML or JSON depending on the Accept header
func sendResponse(ctx *fiber.Ctx, message string) error {
	if ctx.Accepts("application/xml") != "" {
		ctx.Set("Content-Type", "application/xml")
		return ctx.SendString(message)
	}

	if ctx.Accepts("application/json", "json") != "" {
		return ctx.SendString(message)
	}

	// default to xml
	ctx.Set("Content-Type", "application/xml")
	return ctx.SendString(message)
}

// translateMessage translates the timeRange to french
func translateMessage(timeRange string) string {
	switch timeRange {
	case "day":
		return "Au cours de la dernière journée"
	case "week":
		return "La semaine dernière"
	case "month":
		return "Le mois dernier"
	default:
		return "jour"
	}
}

// CreateUser creates a new user in the database
func (a *API) CreateUser(c *fiber.Ctx) error {
	data := new(mdb.User)
	if err := c.BodyParser(data); err != nil {
		return err
	}

	lang := c.Query("lang", "english")

	log.Println("Creating user")
	oid, err := a.mongo.CreateUser(c.Context(), *data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.Status(409)
			msg := buildResponse(c, "User already exists", -1, lang)
			return sendResponse(c, msg)
		}

		c.Status(500)
		msg := buildResponse(c, "Error creating user", -1, lang)
		return sendResponse(c, msg)
	}

	c.Status(201)
	msg := buildResponse(c, fmt.Sprintf("Created user with ID: %s", oid.Hex()), -1, lang)
	return sendResponse(c, msg)
}

// GetRainfall returns the amount of rain reported in either the past day, week, or month
func (a *API) GetRainfallAmount(c *fiber.Ctx) error {
	timeRange := c.Params("timeRange", "day")
	location := c.Query("location", "")
	lang := c.Query("lang", "english")
	amount, err := a.calculateRainfall(c.Context(), timeRange, location)
	if err != nil {
		c.Status(400)
		msg := buildResponse(c, fmt.Sprintf("Error getting %sly rainfall.", timeRange), -1, lang)
		sendResponse(c, msg)
	}

	var message string
	if lang == "french" {
		message = fmt.Sprintf("%s, il a plu %d millimètres.", translateMessage(timeRange), amount)
	} else {
		message = fmt.Sprintf("In the past %s it rained %d millimeters.", timeRange, amount)
	}
	msg := buildResponse(c, message, amount, lang)

	return sendResponse(c, msg)
}

// ReportRain allows a user to report rainfall in mm
func (a *API) ReportRain(c *fiber.Ctx) error {
	data := new(mdb.Rainfall)
	if err := c.BodyParser(data); err != nil {
		return err
	}

	lang := c.Query("lang", "english")

	if a.config.UsePhoneAuth {
		// Check if there is a phone number, if not return an error
		if data.PhoneNumber == "" {
			log.Println("no phone number provided")
			c.Status(400)
			msg := buildResponse(c, "No phone number provided, please call back and try again.", -1, lang)
			return sendResponse(c, msg)
		}

		// Check if the user exists in the DB
		if !a.mongo.CheckUserExists(c.Context(), data.PhoneNumber) {
			log.Println("Attempted use by unauthorized user")
			c.Status(403)
			msg := buildResponse(c, "You are not authorized to use this service.", -1, lang)
			return sendResponse(c, msg)
		}
		
		// Check if the user has already reported today
		if a.mongo.CheckUserReportedToday(c.Context(), data.PhoneNumber) {
			c.Status(429)
			msg := buildResponse(c, "Sorry, You have reached the maximum number of reports for today.", -1, lang)
			return sendResponse(c, msg)
		}
	}

	// Check if the amount is negative
	if data.Amount < 0 {
		log.Println("negative amount")
		c.Status(400)
		msg := buildResponse(c, "Amount cannot be negative", -1, lang)
		return sendResponse(c, msg)
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
	var msg string
	if lang == "french" {
		msg = buildResponse(c, "Merci pour votre rapport!", -1, lang)
	} else {
		msg = buildResponse(c, "Thank you for your report!", -1, lang)
	}
	return sendResponse(c, msg)
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
	} else if len(results) == 0 {
		return 0, nil
	}

	// Iterate over the results and sum the amount of rain for the day
	var rain int
	for _, result := range results {
		rain += result.Amount
	}

	return rain / len(results), nil
}
