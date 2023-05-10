package mdb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rainfall struct {
	ID primitive.ObjectID						 `json:"_id,omitempty" xml:"_id,omitempty" bson:"_id,omitempty"`
	Amount int											 `json:"amount" xml:"amount" bson:"amount"`
	Location string									 `json:"location" xml:"location" bson:"location"`
	ReportedAt primitive.DateTime 	 `json:"reportedAt" xml:"reportedAt" bson:"reportedAt"`
	// The phone number of the user who reported the rainfall.
	// Used to validate one report per user per day.
	PhoneNumber string							 `json:"phoneNumber" xml:"phoneNumber" bson:"phoneNumber"`
}

// Representation of a user in the database.
// Existing within the database is a 
type User struct {
	ID primitive.ObjectID						 `json:"_id,omitempty" xml:"_id,omitempty" bson:"_id,omitempty"`
	PhoneNumber string							 `json:"phoneNumber" xml:"phoneNumber" bson:"phoneNumber"`
	Name string											 `json:"name,omitempty" xml:"name,omitempty" bson:"name,omitempty"`
}