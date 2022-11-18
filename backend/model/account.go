package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

var DefaultRootScopes = []string{"root:*:*"}
var DefaultAdminScopes = []string{"admin:*:*"}
var DefaultUserScopes = []string{"user:*:*"}

type Account struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password,omitempty" json:"-"`
	Scopes      []string      `bson:"scopes" json:"scopes"`
	CreatedTime time.Time     `bson:"createdTime" json:"createdTime"`
	UpdatedTime time.Time     `bson:"updatedTime" json:"updatedTime"`
	UserID      string        `bson:"userID" json:"userID"`
	DeviceIDs   []string      `bson:"deviceIDs" json:"deviceIDs"`
}
