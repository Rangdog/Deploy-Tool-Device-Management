package repository

import "BE_Manage_device/internal/domain/entity"

type MaintenanceSchedulesRepository interface {
	Create(*entity.MaintenanceSchedules) (*entity.MaintenanceSchedules, error)
	GetAllMaintenanceSchedulesByAssetId(assetId int64) ([]*entity.MaintenanceSchedules, error)
}
