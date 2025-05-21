package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLUserRBACRepository struct {
	db *gorm.DB
}

func NewPostgreSQLUserRBACRepository(db *gorm.DB) repository.UserRBACRepository {
	return &PostgreSQLUserRBACRepository{db: db}
}

func (r *PostgreSQLUserRBACRepository) Create(userRBAC *entity.UserRbac, tx *gorm.DB) error {
	if err := tx.Create(userRBAC).Error; err != nil {
		return err
	}
	return nil
}
