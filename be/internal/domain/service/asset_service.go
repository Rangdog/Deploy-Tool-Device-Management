package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"time"
)

type AssetsService struct {
	repo                repository.AssetsRepository
	assertLogrepository repository.AssetsLogRepository
}

func NewAssetsService(repo repository.AssetsRepository, assertLogrepository repository.AssetsLogRepository) *AssetsService {
	return &AssetsService{repo: repo, assertLogrepository: assertLogrepository}
}

func (service *AssetsService) Create(userId int64, assetName string, purchaseDate time.Time, cost float64, owner *int64, warrantExpiry time.Time, status string, serialNumber string, imageURL *string, fileAttachment *string, categoryId int64, departmentId *int64) (*entity.Assets, error) {
	asset := entity.Assets{
		AssetName:      assetName,
		PurchaseDate:   purchaseDate,
		Cost:           cost,
		Owner:          owner,
		WarrantExpiry:  warrantExpiry,
		Status:         status,
		SerialNumber:   serialNumber,
		ImageUpload:    imageURL,
		FileAttachment: fileAttachment,
		CategoryId:     categoryId,
		DepartmentId:   departmentId,
	}
	assetCreate, err := service.repo.Create(&asset)
	if err != nil {
		return nil, err
	}
	assetLog := entity.AssetLog{
		Action:        "Create",
		Timestamp:     time.Now(),
		UserId:        *owner,
		Asset_id:      assetCreate.Id,
		ChangeSummary: "Create",
	}

	_, err = service.assertLogrepository.Create(&assetLog)
	if err != nil {
		return nil, err
	}
	return assetCreate, nil
}
