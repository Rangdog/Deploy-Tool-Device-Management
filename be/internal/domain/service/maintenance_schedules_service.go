package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"errors"
	"time"
)

type MaintenanceSchedulesService struct {
	repo      repository.MaintenanceSchedulesRepository
	assetRepo repository.AssetsRepository
}

func NewMaintenanceSchedulesService(repo repository.MaintenanceSchedulesRepository, assetRepo repository.AssetsRepository) *MaintenanceSchedulesService {
	return &MaintenanceSchedulesService{repo: repo, assetRepo: assetRepo}
}

func (service *MaintenanceSchedulesService) Create(userId int64, assetId int64, startDate, endDate time.Time) (*entity.MaintenanceSchedules, error) {
	assetCheck, err := service.assetRepo.GetAssetById(assetId)
	if err != nil {
		return nil, err
	}
	if assetCheck.Status == "Disposed" || assetCheck.Status == "Retired" || assetCheck.Status == "Under Maintenance" {
		return nil, errors.New("can't set maintenance schedules because status")
	}
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

func (service *MaintenanceSchedulesService) Update(id int64, startDate time.Time, endDate time.Time) (*entity.MaintenanceSchedules, error) {
	maintenace, err := service.repo.Update(id, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return maintenace, nil
}

func (service *MaintenanceSchedulesService) Delete(id int64) error {
	maintenanceCheck, err := service.repo.GetMaintenanceSchedulesById(id)
	if err != nil {
		return err
	}
	if maintenanceCheck.StartDate.After(time.Now()) {
		return errors.New("start date <= now")
	}
	err = service.repo.Delete(id)
	return err
}

func (service *MaintenanceSchedulesService) GetAllMaintenanceSchedules() ([]*entity.MaintenanceSchedules, error) {
	maintenances, err := service.repo.GetAllMaintenanceSchedules()
	if err != nil {
		return nil, err
	}
	return maintenances, nil
}
