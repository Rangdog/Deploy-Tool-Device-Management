package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type AssignmentRepository interface {
	Create(*entity.Assignments) (*entity.Assignments, error)
	Update(id int64, userId, assetId *int64, AssignBy int64) (*entity.Assignments, error)
	GetDB() *gorm.DB
}
