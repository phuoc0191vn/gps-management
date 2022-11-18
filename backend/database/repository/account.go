package repository

import "ctigroupjsc.com/phuocnn/gps-management/model"

type AccountRepository interface {
	All() ([]model.Account, error)
	FindByEmail(email string) (*model.Account, error)

	Save(account model.Account) error

	UpdateByUsername(username string, account model.Account) error

	RemoveByUsername(username string) error
}
