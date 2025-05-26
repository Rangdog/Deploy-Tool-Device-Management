package repository

import "BE_Manage_device/internal/domain/entity"

type MaintenanceNotificationsRepository interface {
	Create(*entity.MaintenanceNotifications) (*entity.MaintenanceNotifications, error)
	GetMaintenanceNotificationByScheduleId(scheduleId int64) (*entity.MaintenanceNotifications, error)
}
