package account

import (
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
)

type DetailAccount struct {
	Name        string `bson:"name" json:"name"`
	IdentityID  string `bson:"identityID" json:"identityID"`
	Unit        string `bson:"unit" json:"unit"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`

	Email       string    `bson:"email" json:"email"`
	CreatedTime time.Time `bson:"createdTime" json:"createdTime"`
	UpdatedTime time.Time `bson:"updatedTime" json:"updatedTime"`
	Scope       string    `bson:"scope" json:"scope"`
	DeviceIDs   []string  `bson:"deviceIDs" json:"deviceIDs"`
	UserID      string    `bson:"userID" json:"userID"`
	CreatedBy   string    `bson:"createdBy" json:"createdBy"`
}

type DetailAccountHandler struct {
	AccountRepository repository.AccountRepository
	UserRepository    repository.UserRepository
}

func (h *DetailAccountHandler) Handle(id string) (DetailAccount, error) {
	account, err := h.AccountRepository.FindByID(id)
	if err != nil {
		return DetailAccount{}, err
	}

	user, err := h.UserRepository.FindByID(account.UserID)
	if err != nil {
		return DetailAccount{}, err
	}

	return DetailAccount{
		Name:        user.Name,
		IdentityID:  user.IdentityID,
		Unit:        user.Unit,
		PhoneNumber: user.PhoneNumber,
		Email:       account.Email,
		CreatedTime: account.CreatedTime,
		UpdatedTime: account.UpdatedTime,
		Scope:       account.Scope,
		DeviceIDs:   account.DeviceIDs,
		UserID:      account.UserID,
		CreatedBy:   account.CreatedBy,
	}, nil
}
