package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	typeVehicleMotorbike = 0
	typeVehicleCar       = 1

	StatusDisable = 0
	StatusEnable  = 1
)

type Device struct {
	ID              bson.ObjectId `bson:"_id" json:"id"`
	IdentifyID      string        `bson:"identifyID" json:"identifyID"`
	Name            string        `bson:"name" json:"name"`
	Type            int           `bson:"type" json:"type"`
	Status          int           `bson:"status" json:"status"`
	CreatedTime     time.Time     `bson:"createdTime" json:"createdTime"`
	UpdatedTime     time.Time     `bson:"updatedTime" json:"updatedTime"`
	AccountID       string        `bson:"accountID" json:"accountID"`
	CurrentLocation Location      `bson:"currentLocation" json:"currentLocation"`
}
