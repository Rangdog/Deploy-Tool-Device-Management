package repository

import "BE_Manage_device/internal/domain/entity"

type NotificationRepository interface {
	Create(*entity.Notifications) (*entity.Notifications, error)
	GetNotificationsByUserId(userId int64) ([]*entity.Notifications, error)
	UpdateStatus(id int64) error
}
