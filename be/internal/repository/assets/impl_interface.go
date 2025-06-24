package repository

import (
	"BE_Manage_device/internal/domain/entity"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostgreSQLAssetsRepository struct {
	db *gorm.DB
}

func NewPostgreSQLAssetsRepository(db *gorm.DB) AssetsRepository {
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

func (r *PostgreSQLAssetsRepository) UpdateAssetLifeCycleStage(id int64, status string, tx *gorm.DB) (*entity.Assets, error) {
	result := tx.Model(&entity.Assets{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return nil, result.Error
	}

	var asset entity.Assets
	if err := tx.First(&asset, id).Error; err != nil {
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

func (r *PostgreSQLAssetsRepository) GetUserHavePermissionNotifications(id int64) ([]*entity.Users, error) {
	users := []*entity.Users{}
	result := r.db.Model(&entity.Assets{}).
		Joins("JOIN user_rbacs ON user_rbacs.asset_id = assets.id").
		Joins("JOIN users ON users.id = user_rbacs.user_id").
		Joins("JOIN roles ON roles.id = user_rbacs.role_id").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("assets.id = ? AND permissions.slug = ?", id, "notifications").
		Select("users.*").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLAssetsRepository) CheckAssetFinishMaintenance(id int64) (bool, error) {
	var schedule entity.MaintenanceSchedules
	err := r.db.
		Debug().
		Model(&entity.MaintenanceSchedules{}).
		Joins("JOIN assets ON assets.id = maintenance_schedules.asset_id").
		Where("assets.id = ?", id).
		Order("maintenance_schedules.start_date DESC").
		First(&schedule).Error

	if err != nil {
		logrus.Printf("⚠️ Error checking maintenance for asset %d: %v", id, err)
		return false, err
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Chuyển EndDate về đúng timezone
	endDateInLoc := schedule.EndDate.In(loc)

	return endDateInLoc.Before(endOfDay), nil
}

func (r *PostgreSQLAssetsRepository) GetAssetByStatus(status string) ([]*entity.Assets, error) {
	var assets = []*entity.Assets{}
	result := r.db.Model(entity.Assets{}).Where("status = ?", status).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

func (r *PostgreSQLAssetsRepository) GetAssetsWasWarrantyExpiry() ([]*entity.Assets, error) {
	assets := []*entity.Assets{}
	loc, _ := time.LoadLocation("Asia/Bangkok")

	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Model(&entity.Assets{}).
		Joins("LEFT JOIN notifications ON notifications.asset_id = assets.id AND notifications.type = ?", "Expired").
		Where("assets.warrant_expiry < ?", endOfDay).
		Where("notifications.id IS NULL").Preload("OnwerUser").
		Find(&assets).Error

	return assets, err
}

func (r *PostgreSQLAssetsRepository) UpdateOwner(id int64, ownerId int64, tx *gorm.DB) error {
	result := tx.Model(entity.Assets{}).Where("id = ?", id).Update("owner", ownerId)
	return result.Error
}

func (r *PostgreSQLAssetsRepository) UpdateAssetDepartment(id, departmentId int64, tx *gorm.DB) (*entity.Assets, error) {
	result := tx.Model(entity.Assets{}).Where("id = ?", id).Update("department_id", departmentId)
	if result.Error != nil {
		return nil, result.Error
	}
	var asset = entity.Assets{}
	tx.Model(entity.Assets{}).Where("id = ?", id).Find(&asset)
	return &asset, nil
}

func (r *PostgreSQLAssetsRepository) UpdateAssetOwner(id, ownerId int64, tx *gorm.DB) (*entity.Assets, error) {
	result := tx.Model(entity.Assets{}).Where("id = ?", id).Update("owner", ownerId)
	if result.Error != nil {
		return nil, result.Error
	}
	var asset = entity.Assets{}
	tx.Model(entity.Assets{}).Where("id = ?", id).Find(&asset)
	return &asset, nil
}

func (r *PostgreSQLAssetsRepository) GetAssetsByCateOfDepartment(categoryId int64, departmentId int64) ([]*entity.Assets, error) {
	assets := []*entity.Assets{}
	result := r.db.Model(entity.Assets{}).Where("category_id = ? and department_id != ?", categoryId, departmentId).Preload("Category").Preload("Department").Preload("OnwerUser").Preload("Department.Location").Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

func (r *PostgreSQLAssetsRepository) UpdateCost(id int64, cost float64) error {
	result := r.db.Model(entity.Assets{}).Where("id = ?", id).Update("cost", cost)
	return result.Error
}

func (r *PostgreSQLAssetsRepository) UpdateAcquisitionDate(id int64, AcquisitionDate time.Time, tx *gorm.DB) error {
	result := tx.Model(entity.Assets{}).Where("id = ?", id).Update("acquisition_date", AcquisitionDate)
	return result.Error
}

func (r *PostgreSQLAssetsRepository) DeleteOwnerAssetOfOwnerId(ownerId int64) error {
	err := r.db.Model(&entity.Assets{}).
		Where("owner = ?", ownerId).
		Update("owner", nil).Error
	return err
}
func (r *PostgreSQLAssetsRepository) GetAllAssetNotHaveMaintenance() ([]*entity.Assets, error) {
	var assets = []*entity.Assets{}
	today := time.Now()
	result := r.db.Model(&entity.Assets{}).
		Joins("LEFT JOIN maintenance_schedules ON maintenance_schedules.asset_id = assets.id").
		Where("(assets.status = ? OR assets.status = ?)", "In Use", "New").
		Where("maintenance_schedules.start_date < ? OR maintenance_schedules.start_date IS NULL", today).
		Distinct("assets.id, assets.*").
		Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil

}

func (r *PostgreSQLAssetsRepository) GetAllAssetOfDep(depId int64) ([]*entity.Assets, error) {
	assets := []*entity.Assets{}
	result := r.db.Model(entity.Assets{}).Where("department_id = ?", depId).Preload("Category").Preload("Department").Preload("OnwerUser").Preload("Department.Location").Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}
