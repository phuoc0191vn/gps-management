package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	ScopeRoot  = "root"
	ScopeAdmin = "admin"
	ScopeUser  = "user"
)

type Account struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password,omitempty" json:"-"`
	Scope       string        `bson:"scope" json:"scope"`
	CreatedTime time.Time     `bson:"createdTime" json:"createdTime"`
	UpdatedTime time.Time     `bson:"updatedTime" json:"updatedTime"`
	UserID      string        `bson:"userID" json:"userID"`
	DeviceIDs   []string      `bson:"deviceIDs" json:"deviceIDs"`
	CreatedBy   string        `bson:"createdBy" json:"createdBy"`
}
