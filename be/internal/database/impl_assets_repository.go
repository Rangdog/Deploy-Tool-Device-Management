package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"errors"

	"gorm.io/gorm"
)

type PostgreSQLAssetsRepository struct {
	db *gorm.DB
}

func NewPostgreSQLAssetsRepository(db *gorm.DB) repository.AssetsRepository {
	return &PostgreSQLAssetsRepository{db: db}
}

func (r *PostgreSQLAssetsRepository) Create(assets *entity.Assets, tx *gorm.DB) (*entity.Assets, error) {
	result := tx.Create(assets)
	return assets, result.Error
}

func (r *PostgreSQLAssetsRepository) GetAssetById(id int64) (*entity.Assets, error) {
	asset := &entity.Assets{}
	result := r.db.Model(&entity.Assets{}).Where("id = ?", id).Preload("Category").Preload("Department").Preload("OnwerUser").Preload("Department.Location").First(asset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("can't find record this id")
		}
		return nil, result.Error
	}
	return asset, nil
}

func (r *PostgreSQLAssetsRepository) Delete(id int64) error {
	result := r.db.Model(entity.Assets{}).Where("id = ?", id).Delete(entity.Assets{})
	return result.Error
}

func (r *PostgreSQLAssetsRepository) UpdateAssetLifeCycleStage(id int64, status string) (*entity.Assets, error) {
	result := r.db.Model(&entity.Assets{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return nil, result.Error
	}

	var asset entity.Assets
	if err := r.db.First(&asset, id).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *PostgreSQLAssetsRepository) GetAllAsset() ([]*entity.Assets, error) {
	assets := []*entity.Assets{}
	result := r.db.Model(entity.Assets{}).Preload("Category").Preload("Department").Preload("OnwerUser").Preload("Department.Location").Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

func (r *PostgreSQLAssetsRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgreSQLAssetsRepository) GetHeadDepartmentIdByAssetId(id int64) (int64, error) {
	department := entity.Departments{}
	result := r.db.Model(entity.Assets{}).Joins("Join departments on departments.id = department_id").Where("id = ?", id).First(&department)
	if result.Error != nil {
		return -1, result.Error
	}
	return *department.UserId, nil
}

func (r *PostgreSQLAssetsRepository) UpdateAsset(assets *entity.Assets, tx *gorm.DB) (*entity.Assets, error) {
	var assetUpdate = entity.Assets{}
	updates := map[string]interface{}{}
	if assets.AssetName != "" {
		updates["asset_name"] = assets.AssetName
	}
	if !assets.PurchaseDate.IsZero() {
		updates["purchase_date"] = assets.PurchaseDate
	}
	updates["cost"] = assets.Cost
	if assets.Owner != nil {
		updates["owner"] = assets.Owner
	}
	if !assets.WarrantExpiry.IsZero() {
		updates["warrant_expiry"] = assets.WarrantExpiry
	}
	if assets.Status != "" {
		updates["status"] = assets.Status
	}
	if assets.SerialNumber != "" {
		updates["serial_number"] = assets.SerialNumber
	}
	if assets.FileAttachment != nil {
		updates["file_attachment"] = assets.FileAttachment
	}
	if assets.ImageUpload != nil {
		updates["image_upload"] = assets.ImageUpload
	}
	if assets.ScheduleMaintenance != nil {
		updates["schedule_maintenance"] = assets.ScheduleMaintenance
	}
	if assets.CategoryId != 0 {
		updates["category_id"] = assets.CategoryId
	}
	if assets.DepartmentId != 0 {
		updates["DepartmentId"] = assets.DepartmentId
	}
	err := tx.Model(&assetUpdate).Where("id = ?", assets.Id).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	// Trả về bản ghi sau khi cập nhật (tuỳ bạn muốn lấy lại hay không)
	err = tx.First(&assetUpdate, assets.Id).Error
	if err != nil {
		return nil, err
	}

	return &assetUpdate, nil
}

func (r *PostgreSQLAssetsRepository) DeleteAsset(id int64, tx *gorm.DB) error {
	result := tx.Model(entity.Assets{}).Where("id = ?", id).Update("status", "Disposed")
	return result.Error
}

func (r *PostgreSQLAssetsRepository) UpdateQrURL(assetId int64, qrUrl string) error {
	result := r.db.Model(entity.Assets{}).Where("id = ?", assetId).Update("qr_url", qrUrl)
	return result.Error
}
