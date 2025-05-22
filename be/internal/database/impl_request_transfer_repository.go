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
	result := r.db.Model(entity.RequestTransfer{}).Create(requestTransfer)
	if result.Error != nil {
		return nil, result.Error
	}
	return requestTransfer, result.Error
}
