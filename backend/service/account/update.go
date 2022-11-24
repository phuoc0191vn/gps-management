package account

import (
	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
)

type UpdateAccount struct {
	UserID      string `bson:"userID" json:"userID"`
	Name        string `bson:"name" json:"name"`
	IdentityID  string `bson:"identityID" json:"identityID"`
	Unit        string `bson:"unit" json:"unit"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
}

func (c *UpdateAccount) Valid() error {
	return nil
}

type UpdateAccountHandler struct {
	AccountRepository repository.AccountRepository
	UserRepository    repository.UserRepository
}

func (h *UpdateAccountHandler) Handle(a *UpdateAccount) error {
	err := a.Valid()
	if err != nil {
		return err
	}

	user, err := h.UserRepository.FindByID(a.UserID)
	if err != nil {
		return err
	}

	if a.IdentityID != "" {
		user.IdentityID = a.IdentityID
	}

	if a.Name != "" {
		user.Name = a.Name
	}

	if a.Unit != "" {
		user.Unit = a.Unit
	}

	if a.PhoneNumber != "" {
		user.PhoneNumber = a.PhoneNumber
	}

	return h.UserRepository.UpdateByID(a.UserID, *user)
}
