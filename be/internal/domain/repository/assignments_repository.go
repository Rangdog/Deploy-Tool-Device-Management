package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type AssignmentRepository interface {
	Create(*entity.Assignments) (*entity.Assignments, error)
	Update(assignmentId int64, AssignBy, assetId int64, userId, departmentId *int64, tx *gorm.DB) (*entity.Assignments, error)
	GetDB() *gorm.DB
}
