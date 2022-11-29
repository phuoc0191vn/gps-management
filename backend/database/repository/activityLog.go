package repository

import (
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/model"
)

type ActivityLogRepository interface {
	All() ([]model.ActivityLog, error)
	GetInDay(deviceID, accountID string, date time.Time) (*model.ActivityLog, error)

	Save(log model.ActivityLog) error

	UpdateByID(id string, log model.ActivityLog) error

	RemoveByID(id string) error
	RemoveByAccountID(accountID string) error
}
