package mdb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Rainfall struct {
	Amount int					`json:"amount" bson:"amount"`
	Location string			`json:"location" bson:"location"`
	ReportedAt string 	`json:"reported_at" bson:"reported_at"`
}