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
	oid, err := a.dbWriteUser(c.Context(), *data)
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
func (a *API) GetRainfall(c *fiber.Ctx) error {
	timeRange := c.Params("timeRange", "day")
	location := c.Query("location", "")
	amount, err := a.dbListRainfall(c.Context(), timeRange, location)
	if err != nil {
		c.XML("Error getting day rainfall")
	}

	message := fmt.Sprintf("In the past %s it rained %d milliliters", timeRange, amount)

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
		if !a.dbCheckUserExists(c.Context(), data.PhoneNumber) {
			log.Println("Attempted use by unauthorized user")
			c.Status(403)
			return c.SendString(convertToVXML("You are not authorized to use this service."))
		}
		
		// Check if the user has already reported today
		if a.dbCheckUserReportedToday(c.Context(), data.PhoneNumber) {
			c.Status(429)
			return c.SendString(convertToVXML("Sorry, You have reached the maximum number of reports for today."))
		}
	}
		
	// Set the reportedAt field to the current time
	data.ReportedAt = primitive.NewDateTimeFromTime(time.Now())

	// INSERT
	oid, err := a.dbWriteRainfall(c.Context(), *data)
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

func (a *API) dbListRainfall(ctx context.Context, timeRange string, loc string) (int, error) {
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

	results, err := a.dbFind(ctx, filter)
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

// Generic function to find documents in the database
// TODO: This should be moved to the mdb package
// TODO: Make more generic
func (a *API) dbFind(ctx context.Context, filter bson.M) ([]*mdb.Rainfall, error) {
	collection := a.mongo.Database(a.config.DbName).Collection(a.config.DbCollection)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// TODO: allow for other types of data
	var results []*mdb.Rainfall

	// Iterate over the cursor and append the results to the results slice
	for cursor.Next(ctx) {
		var element mdb.Rainfall
		// Decode the document, which returns a Rainfall struct from the db response
		if err := cursor.Decode(&element); err != nil {
			log.Println(err)
			return nil, err
		}

		results = append(results, &element)
	}

	// Close the cursor once finished
	cursor.Close(ctx)

	return results, nil
}

// Generic function to write Rainfall documents to the database
func (a *API) dbWriteRainfall(ctx context.Context, data mdb.Rainfall) (primitive.ObjectID, error) {
	collection := a.mongo.Database(a.config.DbName).Collection(a.config.DbCollection)

	insertResult, err := collection.InsertOne(ctx, data)

	if err != nil {
		log.Println(err)
		return primitive.NilObjectID, err
	}

	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	} else {
		log.Println(err)
		return primitive.NilObjectID, err
	}
}

// Generic function to write User documents to the database
func (a *API) dbWriteUser(ctx context.Context, data mdb.User) (primitive.ObjectID, error) {
	collection := a.mongo.Database(a.config.DbName).Collection(a.config.UserCollection)

	insertResult, err := collection.InsertOne(ctx, data)

	if err != nil {
		log.Println(err)
		return primitive.NilObjectID, err
	}

	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	} else {
		log.Println(err)
		return primitive.NilObjectID, err
	}
}

// Db function to check if a user exists by phone number
func (a *API) dbCheckUserExists(ctx context.Context, phoneNumber string) bool {
	collection := a.mongo.Database(a.config.DbName).Collection(a.config.UserCollection)

	filter := bson.M{"phoneNumber": phoneNumber}

	var result mdb.User
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Println(err)
		return false
	}

	return true
}

// Db function to check if a user has reported today
func (a *API) dbCheckUserReportedToday(ctx context.Context, phoneNumber string) bool {
	collection := a.mongo.Database(a.config.DbName).Collection(a.config.UserCollection)
	
	filter := bson.M{
		"reportedAt": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -1))},
		"phoneNumber": phoneNumber,
	}

	var result mdb.Rainfall
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Println(err)
		return false
	}

	return true
}
