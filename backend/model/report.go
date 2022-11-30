package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	StatusUnprocessedReport = 0
	StatusProcessingReport  = 1
	StatusProcessedReport   = 2
)

type Report struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	AccountID   string        `bson:"accountID" json:"accountID"`
	Name        string        `bson:"name" json:"name"`
	Filename    string        `bson:"filename" json:"filename"`
	Status      int           `bson:"status" json:"status"`
	Type        int           `bson:"type" json:"type"`
	CreatedTime time.Time     `bson:"createdTime" json:"createdTime"`
	Filter      Filter        `bson:"filter" json:"filter"`
}

type Filter struct {
	AccountID string    `bson:"accountID" json:"accountID"`
	DeviceID  string    `bson:"deviceID" json:"deviceID"`
	StartTime time.Time `bson:"startTime" json:"startTime"`
	EndTime   time.Time `bson:"endTime" json:"endTime"`
}
