package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.Users) error
	FindByToken(token string) (*entity.Users, error)
	Update(user *entity.Users) error
	UpdatePassword(user *entity.Users) error
	FindByEmail(email string) (*entity.Users, error)
	FindByUserId(userId int64) (*entity.Users, error)
	DeleteUser(email string) error
	GetDB() *gorm.DB
}
