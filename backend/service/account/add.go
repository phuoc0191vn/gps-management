package account

import (
	"fmt"
	"strings"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	utilities "ctigroupjsc.com/phuocnn/gps-management/uitilities"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type AddAccount struct {
	Name        string `bson:"name" json:"name"`
	IdentityID  string `bson:"identityID" json:"identityID" valid:"required"`
	Unit        string `bson:"unit" json:"unit"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`

	Email     string   `bson:"email" json:"email" valid:"required,email"`
	Password  string   `bson:"password" json:"password" valid:"required"`
	Scopes    []string `bson:"scopes" json:"scopes"`
	DeviceIDs []string `bson:"deviceIDs" json:"deviceIDs"`
	CreatedBy string   `bson:"createdBy" json:"createdBy"`
}

func (c *AddAccount) Valid() error {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)

	if c.PhoneNumber != "" && !utilities.IsPhoneNumber(c.PhoneNumber) {
		return fmt.Errorf("invalid phone number")
	}

	if len(c.Scopes) == 0 {
		return fmt.Errorf("invalid scope")
	}

	_, err := govalidator.ValidateStruct(c)
	return err
}

type AddAccountHandler struct {
	AccountRepository repository.AccountRepository
	UserRepository    repository.UserRepository
}

func (h *AddAccountHandler) Handle(c *AddAccount) error {
	if err := c.Valid(); err != nil {
		return err
	}

	u := model.User{
		ID:          bson.NewObjectId(),
		Name:        c.Name,
		IdentityID:  c.IdentityID,
		Unit:        c.Unit,
		PhoneNumber: c.PhoneNumber,
		CreatedTime: utilities.TimeInUTC(time.Now()),
		UpdatedTime: utilities.TimeInUTC(time.Now()),
	}
	err := h.UserRepository.Save(u)
	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		h.UserRepository.RemoveByID(u.ID.Hex())
		return err
	}

	account := model.Account{
		ID:          bson.NewObjectId(),
		Email:       c.Email,
		Password:    string(password),
		Scopes:      c.Scopes,
		CreatedTime: utilities.TimeInUTC(time.Now()),
		UpdatedTime: utilities.TimeInUTC(time.Now()),
		UserID:      u.ID.Hex(),
		DeviceIDs:   c.DeviceIDs,
		CreatedBy:   c.CreatedBy,
	}

	err = h.AccountRepository.Save(account)
	if err != nil {
		h.UserRepository.RemoveByID(u.ID.Hex())
		return err
	}

	return nil
}
