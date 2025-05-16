package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type UsersSessionRepository interface {
	Create(usersSessions *entity.UsersSesions, tx *gorm.DB) error
	FindByRefreshToken(refreshToken string) (*entity.UsersSesions, error)
	UpdateIsRevoked(user *entity.UsersSesions) error
	CheckUserInSession(userId int64) bool
}
