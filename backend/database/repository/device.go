package repository

import "ctigroupjsc.com/phuocnn/gps-management/model"

type DeviceRepository interface {
	All() ([]model.Device, error)

	Save(device model.Device) error

	UpdateByIdentifyID(identifyID string, device model.Device) error

	RemoveByIdentifyID(identifyID string) error
}
