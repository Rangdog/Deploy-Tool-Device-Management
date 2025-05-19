package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type UsersSessionRepository interface {
	Create(usersSessions *entity.UsersSessions, tx *gorm.DB) error
	FindByRefreshToken(refreshToken string) (*entity.UsersSessions, error)
	UpdateIsRevoked(user *entity.UsersSessions) error
	CheckUserInSession(userId int64) bool
	FindByUserIdInSession(UserId int64) (*entity.UsersSessions, error)
	CheckTokenWasInVoked(accessToken string) bool
}
