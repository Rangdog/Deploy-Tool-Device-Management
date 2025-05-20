package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/pkg/utils"
	"fmt"
	"mime/multipart"
	"time"
)

type AssetsService struct {
	repo                repository.AssetsRepository
	assertLogRepository repository.AssetsLogRepository
	roleRepository      repository.RoleRepository
}

func NewAssetsService(repo repository.AssetsRepository, assertLogRepository repository.AssetsLogRepository, roleRepository repository.RoleRepository) *AssetsService {
	return &AssetsService{repo: repo, assertLogRepository: assertLogRepository, roleRepository: roleRepository}
}

func (service *AssetsService) Create(userId int64, assetName string, purchaseDate time.Time, cost float64, owner *int64, warrantExpiry time.Time, status string, serialNumber string, image *multipart.FileHeader, fileAttachment *multipart.FileHeader, categoryId int64, departmentId *int64) (*entity.Assets, error) {
	imgFile, err := image.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open image: %w", err)
	}
	defer imgFile.Close()
	fileFile, err := fileAttachment.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open fileAttachment: %w", err)
	}
	defer fileFile.Close()
	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), image.Filename)
	imagePath := "images/" + uniqueName
	uniqueName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileAttachment.Filename)
	filePath := "files/" + uniqueName
	uploader := utils.NewSupabaseUploader()
	imageUrl, err := uploader.Upload(imagePath, imgFile, image.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	fileUrl, err := uploader.Upload(filePath, fileFile, fileAttachment.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	asset := entity.Assets{
		AssetName:      assetName,
		PurchaseDate:   purchaseDate,
		Cost:           cost,
		Owner:          owner,
		WarrantExpiry:  warrantExpiry,
		Status:         status,
		SerialNumber:   serialNumber,
		ImageUpload:    &imageUrl,
		FileAttachment: &fileUrl,
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
	_, err = service.assertLogRepository.Create(&assetLog)
	if err != nil {
		return nil, err
	}

	return assetCreate, nil
}

func (service *AssetsService) GetAssetById(userId int64, assertId int64) (*entity.Assets, error) {
	assert, err := service.repo.GetAssetById(assertId)
	if err != nil {
		return nil, err
	}
	return assert, err
}

func (service *AssetsService) GetAllAsset() ([]*entity.Assets, error) {
	assets, err := service.repo.GetAllAsset()
	if err != nil {
		return nil, err
	}
	return assets, err
}
