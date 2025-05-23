package service

import (
	"BE_Manage_device/internal/domain/repository"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}
