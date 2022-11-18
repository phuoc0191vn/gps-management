package account

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	utilities "ctigroupjsc.com/phuocnn/gps-management/uitilities"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type AddAccount struct {
	Email     string   `bson:"email" json:"email" valid:"required,email"`
	Password  string   `bson:"password" json:"password" valid:"required"`
	Scopes    []string `bson:"scopes" json:"scopes"`
	UserID    string   `bson:"userID" json:"userID"`
	DeviceIDs []string `bson:"deviceIDs" json:"deviceIDs"`
}

func (c *AddAccount) Valid() error {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)

	if len(c.Scopes) == 0 {
		return fmt.Errorf("invalid scope")
	}

	_, err := govalidator.ValidateStruct(c)
	return err
}

type AddAccountHandler struct {
	AccountRepository repository.AccountRepository
}

func (h *AddAccountHandler) Handle(c *AddAccount) error {
	if err := c.Valid(); err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	account := model.Account{
		ID:          bson.NewObjectId(),
		Email:       c.Email,
		Password:    string(password),
		Scopes:      c.Scopes,
		CreatedTime: utilities.TimeInUTC(time.Now()),
		UpdatedTime: utilities.TimeInUTC(time.Now()),
		UserID:      c.UserID,
		DeviceIDs:   c.DeviceIDs,
	}

	return h.AccountRepository.Save(account)
}
