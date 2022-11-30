package model

import "gopkg.in/mgo.v2/bson"

const (
	StatusContactUnprocessed = 0
	StatusContactProcessed   = 1
)

type Contact struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Email      string        `bson:"email" json:"email"`
	AccountID  string        `bson:"accountID" json:"accountID"`
	AdminEmail string        `bson:"adminEmail" json:"adminEmail"`
	Status     int           `bson:"status" json:"status"`
	Title      string        `bson:"title" json:"title"`
	Body       string        `bson:"body" json:"body"`
}
