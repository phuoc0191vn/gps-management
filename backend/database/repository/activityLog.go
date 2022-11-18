package repository

import "ctigroupjsc.com/phuocnn/gps-management/model"

type ActivityLogRepository interface {
	All() ([]model.ActivityLog, error)

	Save(log model.ActivityLog) error

	RemoveByID(id string) error
}
