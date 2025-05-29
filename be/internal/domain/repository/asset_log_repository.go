package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type AssetsLogRepository interface {
	Create(assetsLog *entity.AssetLog, tx *gorm.DB) (*entity.AssetLog, error)
	GetLogByAssetId(assetId int64) ([]*entity.AssetLog, error)
	GetDB() *gorm.DB
	GetNewLogByAssetId(assetId int64) (*entity.AssetLog, error)
}
