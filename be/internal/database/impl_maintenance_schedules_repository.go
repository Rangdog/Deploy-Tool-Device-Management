package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

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
	result := r.db.Model(entity.MaintenanceSchedules{}).Where("asset_id = ?", assetId).Find(maintenances)
	if result.Error != nil {
		return nil, result.Error
	}
	return maintenances, result.Error
}
