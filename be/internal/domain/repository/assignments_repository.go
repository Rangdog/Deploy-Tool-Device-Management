package repository

import "BE_Manage_device/internal/domain/entity"

type AssignmentRepository interface {
	Create(*entity.Assignments) (*entity.Assignments, error)
	Update(id int64, userId, assetId *int64, assetBy int64) (*entity.Assignments, error)
}
