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

func (r *PostgreSQLAssetsRepository) Create(assets *entity.Assets) (*entity.Assets, error) {
	result := r.db.Create(assets)
	return assets, result.Error
}

func (r *PostgreSQLAssetsRepository) GetAssetById(id int64) (*entity.Assets, error) {
	asset := &entity.Assets{}
	result := r.db.Model(&entity.Assets{}).Where("id = ?", id).First(asset)
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
