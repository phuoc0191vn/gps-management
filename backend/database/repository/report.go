package repository

import "ctigroupjsc.com/phuocnn/gps-management/model"

type ReportRepository interface {
	All() ([]model.Report, error)

	Save(report model.Report) error

	RemoveByID(id string) error
}
