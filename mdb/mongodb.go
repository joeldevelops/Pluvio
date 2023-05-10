package mdb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to MongoDB
func ConnectMongo(url string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}

// Create a unique key Index for the User collection. Idempotent.
func CreateUserIndex(client *mongo.Client, dbName string, collectionName string) error {
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"phoneNumber": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	// Index creation
	_, err := client.Database(dbName).Collection(collectionName).Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	return nil
}