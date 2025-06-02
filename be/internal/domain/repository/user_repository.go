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
	FindByEmailForLogin(email string) (*entity.Users, error)
	FindByUserId(userId int64) (*entity.Users, error)
	DeleteUser(email string) error
	GetDB() *gorm.DB
	GetAllUser() []*entity.Users
	UpdateUser(user *entity.Users) (*entity.Users, error)
	GetUserHeadDepartment(departmentId int64) (*entity.Users, error)
	GetUserAssetManageOfDepartment(departmentId int64) (*entity.Users, error)
	GetAllUserOfDepartment(departmentTd int64) ([]*entity.Users, error)
	UpdateDepartment(userId int64, departmentId int64) error
}
