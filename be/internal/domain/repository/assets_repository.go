package repository

import (
	"BE_Manage_device/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type AssetsRepository interface {
	Create(assets *entity.Assets, tx *gorm.DB) (*entity.Assets, error)
	GetAssetById(id int64) (*entity.Assets, error)
	Delete(id int64) error
	UpdateAssetLifeCycleStage(id int64, status string, tx *gorm.DB) (*entity.Assets, error)
	GetAllAsset() ([]*entity.Assets, error)
	GetDB() *gorm.DB
	UpdateAsset(asset *entity.Assets, tx *gorm.DB) (*entity.Assets, error)
	DeleteAsset(id int64, tx *gorm.DB) error
	UpdateQrURL(assetId int64, qrUrl string) error
	GetUserHavePermissionNotifications(id int64) ([]*entity.Users, error)
	CheckAssetFinishMaintenance(id int64) (bool, error)
	GetAssetByStatus(string) ([]*entity.Assets, error)
	GetAssetsWasWarrantyExpiry() ([]*entity.Assets, error)
	UpdateOwner(id int64, ownerId int64, tx *gorm.DB) error
	UpdateAssetDepartment(id, departmentId int64, tx *gorm.DB) (*entity.Assets, error)
	UpdateAssetOwner(id, Owner int64, tx *gorm.DB) (*entity.Assets, error)
	GetAssetsByCateOfDepartment(categoryId int64, departmentId int64) ([]*entity.Assets, error)
	UpdateCost(id int64, cost float64) error
	UpdateAcquisitionDate(id int64, AcquisitionDate time.Time, tx *gorm.DB) error
	DeleteOwnerAssetOfOwnerId(ownerId int64) error
	GetAllAssetNotHaveMaintenance() ([]*entity.Assets, error)
}
