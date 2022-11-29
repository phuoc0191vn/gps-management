package device

import (
	"fmt"
	"strconv"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	utilities "ctigroupjsc.com/phuocnn/gps-management/uitilities"

	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

type AddDevice struct {
	Name               string `bson:"name" json:"name"`
	LicensePlateNumber string `bson:"licensePlateNumber" json:"licensePlateNumber"`
	Type               string `bson:"type" json:"type"`
	Status             string `bson:"status" json:"status"`
	AccountID          string `bson:"accountID" json:"accountID" valid:"required"`
}

func (c *AddDevice) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

type AddDeviceHandler struct {
	DeviceRepository  repository.DeviceRepository
	AccountRepository repository.AccountRepository
}

// todo: current logic: 1 account only 1 device --> add 1 device will disable all other devices

func (h *AddDeviceHandler) Handle(d *AddDevice) error {
	if err := d.Valid(); err != nil {
		return err
	}

	deviceIDs, err := h.AccountRepository.GetDeviceIDsByID(d.AccountID)
	if err != nil {
		return err
	}

	typeInt, err := strconv.Atoi(d.Type)
	if err != nil || !model.IsValidType(typeInt) {
		return fmt.Errorf("invalid type")
	}
	statusInt, err := strconv.Atoi(d.Status)
	if err != nil || !model.IsValidStatus(statusInt) {
		return fmt.Errorf("invalid status")
	}

	device := model.Device{
		ID:                 bson.NewObjectId(),
		Name:               d.Name,
		LicensePlateNumber: d.LicensePlateNumber,
		Type:               typeInt,
		Status:             statusInt,
		CreatedTime:        utilities.TimeInUTC(time.Now()),
		UpdatedTime:        utilities.TimeInUTC(time.Now()),
		AccountID:          d.AccountID,
		CurrentLocation:    model.Location{},
	}

	err = h.DeviceRepository.Save(device)
	if err != nil {
		return err
	}

	if device.Status == model.StatusEnable {
		for i := 0; i < len(deviceIDs); i++ {
			h.DeviceRepository.UpdateStatus(deviceIDs[i], model.StatusDisable)
		}
	}

	deviceIDs = append(deviceIDs, device.ID.Hex())
	err = h.AccountRepository.UpdateDeviceIDs(d.AccountID, deviceIDs)
	if err != nil {
		h.DeviceRepository.RemoveByID(device.ID.Hex())
		return err
	}

	return nil
}
