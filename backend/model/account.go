package model

import (
	"time"

	utilities "ctigroupjsc.com/phuocnn/gps-management/uitilities"

	"gopkg.in/mgo.v2/bson"
)

const (
	ScopeRoot  = "Root"
	ScopeAdmin = "Admin"
	ScopeUser  = "User"
)

var (
	DefaultRootScopes  = []string{"root:*:*"}
	DefaultAdminScopes = []string{"admin:*:*"}
	DefaultUserScopes  = []string{"user:*:*"}
)

type Account struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password,omitempty" json:"-"`
	Scopes      []string      `bson:"scopes" json:"scopes"`
	CreatedTime time.Time     `bson:"createdTime" json:"createdTime"`
	UpdatedTime time.Time     `bson:"updatedTime" json:"updatedTime"`
	UserID      string        `bson:"userID" json:"userID"`
	DeviceIDs   []string      `bson:"deviceIDs" json:"deviceIDs"`
	CreatedBy   string        `bson:"createdBy" json:"createdBy"`
}

func ConvertScopes(scopes []string) []string {
	if utilities.EqualStringArray(scopes, DefaultRootScopes) {
		return []string{ScopeRoot}
	}

	if utilities.EqualStringArray(scopes, DefaultAdminScopes) {
		return []string{ScopeAdmin}
	}

	if utilities.EqualStringArray(scopes, DefaultUserScopes) {
		return []string{ScopeUser}
	}

	return scopes
}
