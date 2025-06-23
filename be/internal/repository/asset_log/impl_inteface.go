package repository

import (
	"BE_Manage_device/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type PostgreSQLAssetsLogrepository struct {
	db *gorm.DB
}

func NewPostgreSQLAssetsLogRepository(db *gorm.DB) AssetsLogRepository {
	return &PostgreSQLAssetsLogrepository{db: db}
}

func (r *PostgreSQLAssetsLogrepository) Create(assetsLog *entity.AssetLog, tx *gorm.DB) (*entity.AssetLog, error) {
	loc, _ := time.LoadLocation("Asia/Bangkok") // GMT+7
	assetsLog.Timestamp = assetsLog.Timestamp.In(loc)
	result := tx.Create(assetsLog)
	return assetsLog, result.Error
}

func (r *PostgreSQLAssetsLogrepository) GetLogByAssetId(assetId int64) ([]*entity.AssetLog, error) {
	assetLogs := []*entity.AssetLog{}
	result := r.db.Model(entity.AssetLog{}).Where("asset_id = ?", assetId).Find(&assetLogs)
	if result.Error != nil {
		return nil, result.Error
	}
	return assetLogs, nil
}

func (r *PostgreSQLAssetsLogrepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgreSQLAssetsLogrepository) GetNewLogByAssetId(assetId int64) (*entity.AssetLog, error) {
	var assetLogs entity.AssetLog
	result := r.db.Model(entity.AssetLog{}).Where("asset_id = ?", assetId).Order("timestamp DESC").First(&assetLogs)
	if result.Error != nil {
		return nil, result.Error
	}
	return &assetLogs, nil
}
