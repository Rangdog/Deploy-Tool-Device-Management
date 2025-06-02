package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLRequestTransferRepository struct {
	db *gorm.DB
}

func NewPostgreSQLRequestTransferRepository(db *gorm.DB) repository.RequestTransferRepository {
	return &PostgreSQLRequestTransferRepository{db: db}
}

func (r *PostgreSQLRequestTransferRepository) Create(requestTransfer *entity.RequestTransfer) (*entity.RequestTransfer, error) {
	result := r.db.Model(entity.RequestTransfer{}).Preload("User").Preload("User.Department").Preload("Category").Create(requestTransfer)
	if result.Error != nil {
		return nil, result.Error
	}
	request := entity.RequestTransfer{}
	result = r.db.Model(entity.RequestTransfer{}).Where("id = ?", requestTransfer.Id).Preload("User").Preload("User.Department").Preload("Category").First(&request)
	return &request, result.Error
}

func (r *PostgreSQLRequestTransferRepository) UpdateStatusConfirm(id int64, tx *gorm.DB) (*entity.RequestTransfer, error) {
	request := entity.RequestTransfer{}
	result := tx.Model(entity.RequestTransfer{}).Where("id = ?", id).Update("status", "Confirm")
	if result.Error != nil {
		return nil, result.Error
	}
	result = tx.Model(entity.RequestTransfer{}).Where("id = ?", id).Preload("User").Preload("User.Department").Preload("Category").First(&request)
	if result.Error != nil {
		return nil, result.Error
	}
	return &request, nil
}

func (r PostgreSQLRequestTransferRepository) UpdateStatusDeny(id int64, tx *gorm.DB) (*entity.RequestTransfer, error) {
	request := entity.RequestTransfer{}
	result := r.db.Model(entity.RequestTransfer{}).Where("id = ?", id).Update("status", "Deny")
	if result.Error != nil {
		return nil, result.Error
	}
	result = r.db.Model(entity.RequestTransfer{}).Where("id = ?", id).Preload("User").Preload("User.Department").Preload("Category").First(&request)
	if result.Error != nil {
		return nil, result.Error
	}
	return &request, nil
}

func (r *PostgreSQLRequestTransferRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgreSQLRequestTransferRepository) GetRequestTransferById(id int64) (*entity.RequestTransfer, error) {
	request := entity.RequestTransfer{}
	result := r.db.Model(entity.RequestTransfer{}).Where("id = ?", id).Preload("User").Preload("User.Department").Preload("Category").First(&request)
	if result.Error != nil {
		return nil, result.Error
	}
	return &request, nil
}
