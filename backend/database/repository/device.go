package repository

import "ctigroupjsc.com/phuocnn/gps-management/model"

type DeviceRepository interface {
	All() ([]model.Device, error)
	Pagination(page int, limit int) (int, []model.Device, error)
	FindByID(id string) (*model.Device, error)
	FindByAccountID(accountID string) ([]model.Device, error)

	Save(device model.Device) error

	UpdateByID(id string, device model.Device) error

	RemoveByID(id string) error
}
