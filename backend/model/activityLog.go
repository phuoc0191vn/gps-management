package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type ActivityLog struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	AccountID string        `bson:"accountID" json:"accountID"`
	DeviceID  string        `bson:"deviceID" json:"deviceID"`
	Date      time.Time     `bson:"date" json:"date"`
	Locations []Location    `bson:"locations" json:"locations"`
}
