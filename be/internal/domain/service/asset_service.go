package service

import (
	"BE_Manage_device/internal/domain/repository"
	"time"
)

type AssetsService struct {
	repo repository.AssetsRepository
}

func NewAssetsService(repo repository.AssetsRepository) *AssetsService {
	return &AssetsService{repo: repo}
}

func (service *AssetsService) Create(assetName string, purchaseDate time.Time, cost float64, owner *int64, warrantExpiry time.Time, status string) {

}
