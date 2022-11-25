package device

import (
	"fmt"
	"strconv"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	utilities "ctigroupjsc.com/phuocnn/gps-management/uitilities"

	"github.com/asaskevich/govalidator"
)

type UpdateDevice struct {
	Name               string `bson:"name" json:"name"`
	LicensePlateNumber string `bson:"licensePlateNumber" json:"licensePlateNumber"`
	Type               string `bson:"type" json:"type"`
	Status             string `bson:"status" json:"status"`
	AccountID          string `bson:"accountID" json:"accountID" valid:"required"`
}

func (c *UpdateDevice) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

type UpdateDeviceHandler struct {
	DeviceRepository  repository.DeviceRepository
	AccountRepository repository.AccountRepository
}

func (h *UpdateDeviceHandler) Handle(id string, d *UpdateDevice) error {
	if err := d.Valid(); err != nil {
		return err
	}

	device, err := h.DeviceRepository.FindByID(id)
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

	device.Type = typeInt
	device.Status = statusInt
	if d.Name != "" {
		device.Name = d.Name
	}
	if d.LicensePlateNumber != "" {
		device.LicensePlateNumber = d.LicensePlateNumber
	}

	if d.AccountID == "" || d.AccountID == device.AccountID {
		return h.DeviceRepository.UpdateByID(id, *device)
	}

	// remove deviceID from old-account
	oldAccount, err := h.AccountRepository.GetDeviceIDsByID(device.AccountID)
	if err != nil {
		return err
	}

	tmp := make([]string, 0)
	for i := 0; i < len(oldAccount); i++ {
		if oldAccount[i] == device.ID.Hex() {
			continue
		}

		tmp = append(tmp, oldAccount[i])
	}
	oldAccount = make([]string, 0)
	oldAccount = tmp

	err = h.AccountRepository.UpdateDeviceIDs(device.AccountID, oldAccount)
	if err != nil {
		return err
	}

	// update deviceID into new-account
	device.AccountID = d.AccountID
	newAccount, err := h.AccountRepository.GetDeviceIDsByID(d.AccountID)
	if err != nil {
		return err
	}
	newAccount = append(newAccount, device.ID.Hex())

	err = h.AccountRepository.UpdateDeviceIDs(d.AccountID, newAccount)
	if err != nil {
		return err
	}

	device.UpdatedTime = utilities.TimeInUTC(time.Now())
	return h.DeviceRepository.UpdateByID(id, *device)
}
