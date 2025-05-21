package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
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
