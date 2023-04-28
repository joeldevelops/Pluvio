package mdb

type Rainfall struct {
	Amount int					`json:"amount" bson:"amount"`
	Location string			`json:"location" bson:"location"`
	ReportedAt string 	`json:"reported_at" bson:"reported_at"`
}