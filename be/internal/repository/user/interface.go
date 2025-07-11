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
	GetAllUser(companyId int64) []*entity.Users
	UpdateUser(user *entity.Users) (*entity.Users, error)
	GetUserAssetManageOfDepartment(departmentId int64) (*entity.Users, error)
	GetAllUserRoleEmployeeOfDepartment(departmentTd int64) ([]*entity.Users, error)
	GetAllUserRoleManagerOfDepartment(departmentTd int64) ([]*entity.Users, error)
	UpdateDepartment(userId int64, departmentId int64) error
	CheckHeadDep(depId int64) error
	CheckManagerDep(depId int64) error
	UpdateHeadDep(id int64, isHead bool) error
	UpdateManagerDep(id int64, isManager bool) error
	UpdateCanExport(id int64, canExport bool) error
	GetUserNotHaveDep() ([]*entity.Users, error)
	GetUserRoleAdmin() ([]*entity.Users, error)
	FindManager(userId int64) (*entity.Users, error)
}
