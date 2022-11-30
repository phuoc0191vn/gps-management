package repository

import (
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/model"
)

type ActivityLogRepository interface {
	All() ([]model.ActivityLog, error)
	GetInDay(deviceID, accountID string, date time.Time) (*model.ActivityLog, error)
	GetRange(deviceID, accountID string, start, end time.Time) ([]model.ActivityLog, error)
	FilterReport(filter model.Filter) (chan model.ActivityLogData, chan error)

	Save(log model.ActivityLog) error

	UpdateByID(id string, log model.ActivityLog) error

	RemoveByID(id string) error
	RemoveByAccountID(accountID string) error
}
