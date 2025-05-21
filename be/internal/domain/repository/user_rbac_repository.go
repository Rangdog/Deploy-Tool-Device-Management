package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRBACRepository interface {
	Create(userRBAC *entity.UserRbac, tx *gorm.DB) error
}
