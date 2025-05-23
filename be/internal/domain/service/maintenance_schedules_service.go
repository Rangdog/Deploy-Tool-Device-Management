package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"time"
)

type MaintenanceSchedulesService struct {
	repo repository.MaintenanceSchedulesRepository
}

func NewMaintenanceSchedulesService(repo repository.MaintenanceSchedulesRepository) *MaintenanceSchedulesService {
	return &MaintenanceSchedulesService{repo: repo}
}

func (service *MaintenanceSchedulesService) Create(userId int64, assetId int64, startDate, endDate time.Time) (*entity.MaintenanceSchedules, error) {
	maintenance := entity.MaintenanceSchedules{
		AssetId:   assetId,
		StartDate: startDate,
		EndDate:   endDate,
	}
	maintenanceCreate, err := service.repo.Create(&maintenance)
	if err != nil {
		return nil, err
	}
	return maintenanceCreate, nil
}

func (service *MaintenanceSchedulesService) GetAllMaintenanceSchedulesByAssetId(userId int64, assetId int64) ([]*entity.MaintenanceSchedules, error) {
	maintenances, err := service.repo.GetAllMaintenanceSchedulesByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	return maintenances, nil
}
