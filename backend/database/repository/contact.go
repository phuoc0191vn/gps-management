package repository

import "ctigroupjsc.com/phuocnn/gps-management/model"

type ContactRepository interface {
	All() ([]model.Contact, error)
	Pagination(page int, limit int, condition map[string]interface{}) (int, []model.Contact, error)
	FindByID(id string) (*model.Contact, error)

	Save(contact model.Contact) error

	UpdateStatus(id string, status int) error

	RemoveByID(id string) error
}
