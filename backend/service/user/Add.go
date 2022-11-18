package user

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	utilities "ctigroupjsc.com/phuocnn/gps-management/uitilities"

	"github.com/asaskevich/govalidator"
)

var regexPhoneNumber = regexp.MustCompile(`(?m)(84|0[0-9])+([0-9]{8})\b`)

type AddUser struct {
	Name        string `bson:"name" json:"name"`
	IdentityID  string `bson:"identityID" json:"identityID" valid:"required"`
	Unit        string `bson:"unit" json:"unit"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
}

func (c *AddUser) Valid() error {
	if c.PhoneNumber != "" && !regexPhoneNumber.MatchString(c.PhoneNumber) {
		return fmt.Errorf("invalid phone number")
	}

	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return err
	}

	return nil
}

type AddUserHandler struct {
	UserRepository repository.UserRepository
}

func (h *AddUserHandler) Handle(c *AddUser) (string, error) {
	if err := c.Valid(); err != nil {
		return "", err
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

	return u.ID.Hex(), h.UserRepository.Save(u)
}
