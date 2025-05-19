package repository

import "BE_Manage_device/internal/domain/entity"

type AssetsLogRepository interface {
	Create(*entity.AssetLog) (*entity.AssetLog, error)
}
