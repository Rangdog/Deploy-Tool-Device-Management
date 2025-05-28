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

func (r *PostgreSQLAssignmentRepository) Create(assignment *entity.Assignments, tx *gorm.DB) (*entity.Assignments, error) {
	result := tx.Create(assignment)
	if result.Error != nil {
		return nil, result.Error
	}
	var assignmentCreated = &entity.Assignments{}
	result = tx.Model(entity.Assignments{}).Where("id = ?", assignment.Id).Preload("UserAssigned").Preload("UserAssign").Preload("Asset").Preload("Department").First(&assignmentCreated)
	return assignmentCreated, result.Error
}

func (r *PostgreSQLAssignmentRepository) Update(assignmentId int64, AssignBy, assetId int64, userId, departmentId *int64, tx *gorm.DB) (*entity.Assignments, error) {
	var assignment entity.Assignments

	updates := map[string]interface{}{}
	if userId != nil {
		updates["user_id"] = *userId
	}
	if departmentId != nil {
		updates["department_id"] = *departmentId
	}
	updates["asset_id"] = assetId
	updates["assign_by"] = AssignBy

	err := tx.Model(&assignment).Where("id = ?", assignmentId).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Trả về bản ghi sau khi cập nhật (tuỳ bạn muốn lấy lại hay không)
	err = tx.Preload("UserAssigned").Preload("UserAssign").Preload("Asset").Preload("Department").Preload("Department.Location").First(&assignment, assignmentId).Error
	if err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (r *PostgreSQLAssignmentRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgreSQLAssignmentRepository) GetAssignmentById(id int64) (*entity.Assignments, error) {
	assignment := entity.Assignments{}
	result := r.db.Model(entity.Assignments{}).Where("id = ?", id).Preload("UserAssigned").Preload("UserAssign").Preload("Asset").Preload("Department").Preload("Department.Location").First(&assignment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &assignment, nil
}

func (r *PostgreSQLAssignmentRepository) GetAssignmentByAssetId(assetId int64) (*entity.Assignments, error) {
	assignment := entity.Assignments{}
	result := r.db.Model(entity.Assignments{}).Where("asset_id = ?", assetId).Preload("UserAssigned").Preload("UserAssign").Preload("Asset").Preload("Department").Preload("Department.Location").First(&assignment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &assignment, nil
}
