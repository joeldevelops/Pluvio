package mdb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MDBConfig struct {
	DbName string
	RainCollection string
	UsersCollection string
}

type MongoDB struct {
	*mongo.Client
	config *MDBConfig
}

// Connect to MongoDB
func NewMongoConnection(url string, c *MDBConfig) (*MongoDB, error) {

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return &MongoDB{client, c}, nil
}

// Create a unique key Index for the User collection. Idempotent.
func (db *MongoDB) CreateUserIndex() error {
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"phoneNumber": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	// Index creation
	_, err := db.Database(db.config.DbName).Collection(db.config.UsersCollection).Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	return nil
}

// Get Rainfall documents from the database
func (db *MongoDB) GetRainfall(ctx context.Context, filter bson.M) ([]*Rainfall, error) {
	collection := db.Database(db.config.DbName).Collection(db.config.RainCollection)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// TODO: allow for other types of data
	var results []*Rainfall

	// Iterate over the cursor and append the results to the results slice
	for cursor.Next(ctx) {
		var element Rainfall
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

// Create Rainfall documents in the database
func (db *MongoDB) CreateRainfall(ctx context.Context, data Rainfall) (primitive.ObjectID, error) {
	collection := db.Database(db.config.DbName).Collection(db.config.RainCollection)

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
func (db *MongoDB) CheckUserExists(ctx context.Context, phoneNumber string) bool {
	collection := db.Database(db.config.DbName).Collection(db.config.UsersCollection)

	filter := bson.M{"phoneNumber": phoneNumber}

	var result User
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
func (db *MongoDB) CheckUserReportedToday(ctx context.Context, phoneNumber string) bool {
	collection := db.Database(db.config.DbName).Collection(db.config.RainCollection)
	
	filter := bson.M{
		"reportedAt": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -1))},
		"phoneNumber": phoneNumber,
	}

	var result Rainfall
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

// Get User documents from the database
func (db *MongoDB) CreateUser(ctx context.Context, data User) (primitive.ObjectID, error) {
	collection := db.Database(db.config.DbName).Collection(db.config.UsersCollection)

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