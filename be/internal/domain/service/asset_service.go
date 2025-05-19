package service

import (
	"BE_Manage_device/internal/domain/repository"
)

type AssetsService struct {
	repo repository.AssetsRepository
}

func NewAssetsService(repo repository.AssetsRepository) *AssetsService {
	return &AssetsService{repo: repo}
}

// func (service *AssetsService) Create(assetName string, purchaseDate time.Time, cost float64, owner *int64, warrantExpiry time.Time, status string, serialNumber string, imageURL *string, fileAttachment *string, categoryId int64, DepartmentId *int64) {
// 	asset := entity.Assets{
// 		AssetName:      assetName,
// 		PurchaseDate:   purchaseDate,
// 		Cost:           cost,
// 		Owner:          owner,
// 		WarrantExpiry:  warrantExpiry,
// 		Status:         status,
// 		SerialNumber:   serialNumber,
// 		ImageUpload:    imageURL,
// 		FileAttachment: fileAttachment,
// 	}
// }
