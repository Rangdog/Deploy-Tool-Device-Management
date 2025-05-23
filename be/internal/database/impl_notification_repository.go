package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLNotificationRepository struct {
	db *gorm.DB
}

func NewPostgreSQLNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &PostgreSQLNotificationRepository{db: db}
}

func (r *PostgreSQLNotificationRepository) Create(notification *entity.Notifications) (*entity.Notifications, error) {
	result := r.db.Create(notification)
	if result.Error != nil {
		return nil, result.Error
	}
	return notification, nil
}

func (r *PostgreSQLNotificationRepository) GetNotificationsByUserId(userId int64) ([]*entity.Notifications, error) {
	notifications := []*entity.Notifications{}
	result := r.db.Model(entity.Notifications{}).Where("user_id = ?", userId).Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}
	return notifications, nil
}

func (r *PostgreSQLNotificationRepository) UpdateStatus(id int64) error {
	result := r.db.Model(entity.Notifications{}).Where("id = ?", id).Update("status", "seen")
	return result.Error
}
