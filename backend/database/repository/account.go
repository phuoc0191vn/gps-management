package repository

import "ctigroupjsc.com/phuocnn/gps-management/model"

type AccountRepository interface {
	All() ([]model.Account, error)
	Pagination(page int, limit int, condition map[string]interface{}) (int, []model.Account, error)
	FindByID(id string) (*model.Account, error)
	FindByEmail(email string) (*model.Account, error)

	GetDeviceIDsByID(id string) ([]string, error)
	GetChildAccounts(email string) ([]model.Account, error)

	Save(account model.Account) error

	UpdateByEmail(email string, account model.Account) error
	UpdateDeviceIDs(id string, deviceIDs []string) error

	RemoveByID(id string) error
	RemoveByEmail(email string) error
	RemoveByUserID(userID string) error
}
