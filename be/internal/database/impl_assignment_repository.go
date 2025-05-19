package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLAssignmentRepository struct {
	db *gorm.DB
}

func NewPostgreSQLAssignmentRepository(db *gorm.DB) repository.AssignmentRepository {
	return &PostgreSQLAssignmentRepository{db: db}
}

func (r *PostgreSQLAssignmentRepository) Create(*entity.Assignments) (*entity.Assignments, error)

func (r *PostgreSQLAssignmentRepository) Update(id int64, userId, assetId *int64, assetBy int64) (*entity.Assignments, error)
