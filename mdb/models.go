package mdb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rainfall struct {
	ID primitive.ObjectID						 `json:"_id,omitempty" xml:"_id,omitempty" bson:"_id,omitempty"`
	Amount int											 `json:"amount" xml:"amount" bson:"amount"`
	Location string									 `json:"location" xml:"location" bson:"location"`
	ReportedAt primitive.DateTime 	 `json:"reportedAt" xml:"reportedAt" bson:"reportedAt"`
}