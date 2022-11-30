package repository

import "ctigroupjsc.com/phuocnn/gps-management/model"

type ReportRepository interface {
	All() ([]model.Report, error)
	FindByID(id string) (*model.Report, error)
	Pagination(page int, limit int, condition map[string]interface{}) (int, []model.Report, error)
	GetReportByStatus(status int) ([]model.Report, error)

	Save(report model.Report) error

	UpdateStatusReport(id string, status int) error

	RemoveByID(id string) error
	RemoveByAccountID(accountID string) error
}
