package repository

import (
	"BE_Manage_device/internal/domain/entity"
	"time"
)

type MaintenanceSchedulesRepository interface {
	Create(*entity.MaintenanceSchedules) (*entity.MaintenanceSchedules, error)
	GetAllMaintenanceSchedulesByAssetId(assetId int64) ([]*entity.MaintenanceSchedules, error)
	Update(id int64, startDate time.Time, endDate time.Time) (*entity.MaintenanceSchedules, error)
	Delete(id int64) error
	GetMaintenanceSchedulesById(id int64) (*entity.MaintenanceSchedules, error)
	GetAllMaintenanceSchedules() ([]*entity.MaintenanceSchedules, error)
	GetDateMaintenanceSchedulesInFuture(assetId int64) ([]*entity.TimeRange, error)
}
