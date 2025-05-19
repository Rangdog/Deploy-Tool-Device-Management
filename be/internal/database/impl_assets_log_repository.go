package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLAssetsLogrepository struct {
	db *gorm.DB
}

func NewPostgreSQLAssetsLogRepository(db *gorm.DB) repository.AssetsLogRepository {
	return &PostgreSQLAssetsLogrepository{db: db}
}

func (r *PostgreSQLAssetsLogrepository) Create(assetsLog *entity.AssetLog) (*entity.AssetLog, error) {
	result := r.db.Create(assetsLog)
	return assetsLog, result.Error
}
