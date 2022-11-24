package model

import "gopkg.in/mgo.v2/bson"

const (
	statusGeneratingReport = 0
	statusGeneratedReport  = 1
)

type Report struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	AccountID string        `bson:"accountID" json:"accountID"`
	Name      string        `bson:"name" json:"name"`
	Status    int           `bson:"status" json:"status"`
	Type      int           `bson:"type" json:"type"`
}
