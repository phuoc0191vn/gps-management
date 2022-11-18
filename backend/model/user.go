package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	IdentityID  string        `bson:"identityID" json:"identityID"`
	Unit        string        `bson:"unit" json:"unit"`
	PhoneNumber string        `bson:"phoneNumber" json:"phoneNumber"`
	CreatedTime time.Time     `bson:"createdTime" json:"createdTime"`
	UpdatedTime time.Time     `bson:"updatedTime" json:"updatedTime"`
}
