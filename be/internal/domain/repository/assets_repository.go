package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type AssetsRepository interface {
	Create(assets *entity.Assets, tx *gorm.DB) (*entity.Assets, error)
	GetAssetById(id int64) (*entity.Assets, error)
	Delete(id int64) error
	UpdateAssetLifeCycleStage(id int64, status string, tx *gorm.DB) (*entity.Assets, error)
	GetAllAsset() ([]*entity.Assets, error)
	GetDB() *gorm.DB
	GetHeadDepartmentIdByAssetId(id int64) (int64, error)
	UpdateAsset(asset *entity.Assets, tx *gorm.DB) (*entity.Assets, error)
	DeleteAsset(id int64, tx *gorm.DB) error
	UpdateQrURL(assetId int64, qrUrl string) error
	GetUserHavePermissionNotifications(id int64) ([]*entity.Users, error)
	CheckAssetFinishMaintenance(id int64) (bool, error)
	GetAssetByStatus(string) ([]*entity.Assets, error)
}
