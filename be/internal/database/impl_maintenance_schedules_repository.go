package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type PostgreSQLMaintenanceSchedulesRepository struct {
	db *gorm.DB
}

func NewPostgreSQLMaintenanceSchedulesRepository(db *gorm.DB) repository.MaintenanceSchedulesRepository {
	return &PostgreSQLMaintenanceSchedulesRepository{db: db}
}

func (r *PostgreSQLMaintenanceSchedulesRepository) Create(maintenance *entity.MaintenanceSchedules) (*entity.MaintenanceSchedules, error) {
	result := r.db.Create(maintenance)
	return maintenance, result.Error
}

func (r *PostgreSQLMaintenanceSchedulesRepository) GetAllMaintenanceSchedulesByAssetId(assetId int64) ([]*entity.MaintenanceSchedules, error) {
	maintenances := []*entity.MaintenanceSchedules{}
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	result := r.db.Model(entity.MaintenanceSchedules{}).Where("asset_id = ?", assetId).Where("end_date >= ?", endOfDay).Preload("Asset").Find(&maintenances)
	if result.Error != nil {
		return nil, result.Error
	}
	return maintenances, result.Error
}

func (r *PostgreSQLMaintenanceSchedulesRepository) Update(id int64, startDate time.Time, endDate time.Time) (*entity.MaintenanceSchedules, error) {
	maintenance := entity.MaintenanceSchedules{}
	result := r.db.Model(entity.MaintenanceSchedules{}).Where("id = ?", id).Update("start_date", startDate).Update("end_date", endDate)
	if result.Error != nil {
		return nil, result.Error
	}
	r.db.Model(entity.MaintenanceSchedules{}).Where("id = ?", id).Preload("Asset").First(&maintenance)
	return &maintenance, nil
}

func (r *PostgreSQLMaintenanceSchedulesRepository) Delete(id int64) error {
	if err := r.db.Delete(&entity.MaintenanceSchedules{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgreSQLMaintenanceSchedulesRepository) GetMaintenanceSchedulesById(id int64) (*entity.MaintenanceSchedules, error) {
	maintenance := entity.MaintenanceSchedules{}
	result := r.db.Model(entity.MaintenanceSchedules{}).Where("id = ?", id).Preload("Asset").Preload("Asset.OnwerUser").First(&maintenance)
	if result.Error != nil {
		return nil, result.Error
	}
	return &maintenance, nil
}

func (r *PostgreSQLMaintenanceSchedulesRepository) GetAllMaintenanceSchedules() ([]*entity.MaintenanceSchedules, error) {
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	maintenances := []*entity.MaintenanceSchedules{}
	result := r.db.Model(entity.MaintenanceSchedules{}).Joins("join assets on assets.id = maintenance_schedules.asset_id").Where("end_date >= ?", endOfDay).Preload("Asset").Find(&maintenances)
	if result.Error != nil {
		return nil, result.Error
	}
	return maintenances, result.Error
}

func (r *PostgreSQLMaintenanceSchedulesRepository) GetDateMaintenanceSchedulesInFuture(assetId int64) ([]*entity.TimeRange, error) {
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	maintenances := []*entity.MaintenanceSchedules{}
	result := r.db.Model(entity.MaintenanceSchedules{}).Where("asset_id = ?", assetId).Where("start_date >= ?", endOfDay).Find(&maintenances)
	if result.Error != nil {
		return nil, result.Error
	}
	var TimeRange = []*entity.TimeRange{}
	for _, i := range maintenances {
		TimeRange = append(TimeRange, &entity.TimeRange{Start: i.StartDate, End: i.EndDate})
	}
	return TimeRange, nil
}
