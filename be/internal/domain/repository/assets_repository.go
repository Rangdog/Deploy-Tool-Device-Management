package repository

import "BE_Manage_device/internal/domain/entity"

type AssetsRepository interface {
	Create(*entity.Assets) (*entity.Assets, error)
	GetAssetById(id int64) (*entity.Assets, error)
	Delete(id int64) error
	UpdateAssetLifeCycleStage(id int64, status string) (*entity.Assets, error)
}
