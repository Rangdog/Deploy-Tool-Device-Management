package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type AssetsLogRepository interface {
	Create(assetsLog *entity.AssetLog, tx *gorm.DB) (*entity.AssetLog, error)
}
