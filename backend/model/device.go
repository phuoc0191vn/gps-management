package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	TypeVehicleMotorbike = 0
	TypeVehicleCar       = 1

	StatusDisable = 0
	StatusEnable  = 1
)

type Device struct {
	ID                 bson.ObjectId `bson:"_id" json:"id"`
	Name               string        `bson:"name" json:"name"`
	LicensePlateNumber string        `bson:"licensePlateNumber" json:"licensePlateNumber"`
	Type               int           `bson:"type" json:"type"`
	Status             int           `bson:"status" json:"status"`
	CreatedTime        time.Time     `bson:"createdTime" json:"createdTime"`
	UpdatedTime        time.Time     `bson:"updatedTime" json:"updatedTime"`
	AccountID          string        `bson:"accountID" json:"accountID"`
	CurrentLocation    Location      `bson:"currentLocation" json:"currentLocation"`
}

func IsValidType(typeInt int) bool {
	if typeInt != TypeVehicleCar && typeInt != TypeVehicleMotorbike {
		return false
	}

	return true
}

func IsValidStatus(statusInt int) bool {
	if statusInt != StatusDisable && statusInt != StatusEnable {
		return false
	}

	return true
}
