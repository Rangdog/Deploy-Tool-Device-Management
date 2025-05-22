package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/repository"
	"math"
)

type AssetLogService struct {
	repo repository.AssetsLogRepository
}

func NewAssetLogService(repo repository.AssetsLogRepository) *AssetLogService {
	return &AssetLogService{repo: repo}
}

func (service *AssetLogService) GetLogByAssetId(assetId int64) ([]*entity.AssetLog, error) {
	assetlogs, err := service.repo.GetLogByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	return assetlogs, nil
}
func (service *AssetLogService) Filter(userId int64, assetId int64, action, startTime, endTime *string, page int, limit int) (*map[string]any, error) {
	var filter = filter.AssetLogFilter{
		Action:    action,
		StartTime: startTime,
		EndTime:   endTime,
		Page:      page,
		Limit:     limit,
	}
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.AssetLog{}), assetId)
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	var total int64
	dbFilter.Count(&total)
	offset := (filter.Page - 1) * filter.Limit
	var asset_logs []entity.AssetLog
	resutl := dbFilter.Offset(offset).Limit(filter.Limit).Find(&asset_logs)
	if resutl.Error != nil {
		return nil, resutl.Error
	}

	data := map[string]any{
		"data":       resutl,
		"page":       filter.Page,
		"limit":      filter.Limit,
		"total":      total,
		"total_page": int(math.Ceil(float64(total) / float64(filter.Limit))),
	}
	return &data, nil
}
