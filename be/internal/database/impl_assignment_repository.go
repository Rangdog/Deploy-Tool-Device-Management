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

func (r *PostgreSQLAssignmentRepository) Create(assignment *entity.Assignments) (*entity.Assignments, error) {
	result := r.db.Create(assignment)
	return assignment, result.Error
}

func (r *PostgreSQLAssignmentRepository) Update(id int64, userId, assetId *int64, assetBy int64) (*entity.Assignments, error) {
	var assignment entity.Assignments

	updates := map[string]interface{}{}
	if userId != nil {
		updates["user_id"] = *userId
	}
	if assetId != nil {
		updates["asset_id"] = *assetId
	}
	updates["asset_by"] = assetBy

	err := r.db.Model(&assignment).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Trả về bản ghi sau khi cập nhật (tuỳ bạn muốn lấy lại hay không)
	err = r.db.First(&assignment, id).Error
	if err != nil {
		return nil, err
	}

	return &assignment, nil
}
