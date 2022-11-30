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

type ActivityLogData struct {
	Date      time.Time `bson:"date" json:"date"`
	Longitude float64   `bson:"longitude" json:"longitude"`
	Latitude  float64   `bson:"latitude" json:"latitude"`
	Address   string    `bson:"address" json:"address"`
	TimeStamp int64     `bson:"timestamp" json:"timestamp"`
	Speed     string    `bson:"speed" json:"speed"`
}
