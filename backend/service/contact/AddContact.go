package contact

import (
	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

type AddContact struct {
	Email      string `json:"email" valid:"email"`
	AccountID  string `json:"accountID"`
	AdminEmail string `json:"adminEmail" valid:"email"`
	Status     int    `json:"status"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

func (c *AddContact) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

type AddContactHandler struct {
	ContactRepository repository.ContactRepository
}

func (h *AddContactHandler) Handle(c *AddContact) error {
	if err := c.Valid(); err != nil {
		return err
	}

	contact := model.Contact{
		ID:         bson.NewObjectId(),
		Email:      c.Email,
		AccountID:  c.AccountID,
		AdminEmail: c.AdminEmail,
		Status:     model.StatusContactUnprocessed,
		Title:      c.Title,
		Body:       c.Title,
	}

	return h.ContactRepository.Save(contact)
}
