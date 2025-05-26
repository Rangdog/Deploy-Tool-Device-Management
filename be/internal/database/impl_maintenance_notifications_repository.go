package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLMaintenanceNotificationRepository struct {
	db *gorm.DB
}

func NewPostgreSQLMaintenanceNotificationRepository(db *gorm.DB) repository.MaintenanceNotificationsRepository {
	return &PostgreSQLMaintenanceNotificationRepository{db: db}
}

func (r *PostgreSQLMaintenanceNotificationRepository) Create(maintenanceNotification *entity.MaintenanceNotifications) (*entity.MaintenanceNotifications, error) {
	result := r.db.Create(maintenanceNotification)
	if result.Error != nil {
		return nil, result.Error
	}
	return maintenanceNotification, nil
}

func (r *PostgreSQLMaintenanceNotificationRepository) GetMaintenanceNotificationByScheduleId(scheduleId int64) (*entity.MaintenanceNotifications, error) {
	maintenanceNotification := entity.MaintenanceNotifications{}
	result := r.db.Model(entity.MaintenanceNotifications{}).Where("schedule_id=?").First(&maintenanceNotification)
	if result.Error != nil {
		return nil, result.Error
	}
	return &maintenanceNotification, nil
}
