package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/pkg/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"gorm.io/gorm"
)

type AssetsService struct {
	repo                repository.AssetsRepository
	assertLogRepository repository.AssetsLogRepository
	roleRepository      repository.RoleRepository
	userRBACRepository  repository.UserRBACRepository
	userRepository      repository.UserRepository
}

func NewAssetsService(repo repository.AssetsRepository, assertLogRepository repository.AssetsLogRepository, roleRepository repository.RoleRepository, userRBACRepository repository.UserRBACRepository, userRepository repository.UserRepository) *AssetsService {
	return &AssetsService{repo: repo, assertLogRepository: assertLogRepository, roleRepository: roleRepository, userRBACRepository: userRBACRepository, userRepository: userRepository}
}

func (service *AssetsService) Create(userId int64, assetName string, purchaseDate time.Time, cost float64, owner *int64, warrantExpiry time.Time, serialNumber string, image *multipart.FileHeader, fileAttachment *multipart.FileHeader, categoryId int64, departmentId int64) (*entity.Assets, error) {
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
	tx := service.repo.GetDB().Begin()
	asset := entity.Assets{
		AssetName:      assetName,
		PurchaseDate:   purchaseDate,
		Cost:           cost,
		Owner:          owner,
		WarrantExpiry:  warrantExpiry,
		Status:         "New",
		SerialNumber:   serialNumber,
		ImageUpload:    &imageUrl,
		FileAttachment: &fileUrl,
		CategoryId:     categoryId,
		DepartmentId:   departmentId,
	}
	assetCreate, err := service.repo.Create(&asset, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	assetLog := entity.AssetLog{
		Action:        "Create",
		Timestamp:     time.Now(),
		UserId:        userId,
		Asset_id:      assetCreate.Id,
		ChangeSummary: "Create",
	}
	_, err = service.assertLogRepository.Create(&assetLog, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go service.SetRole(assetCreate.Id, tx, &wg) // commit ở đây
	wg.Wait()
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

func (service *AssetsService) SetRole(assetId int64, tx *gorm.DB, wg *sync.WaitGroup) bool {
	defer wg.Done()
	users := service.userRepository.GetAllUser()
	for _, user := range users {
		if service.roleRepository.GetTitleByRoleId(user.RoleId) == "Department Head" {
			userId, err := service.repo.GetHeadDepartmentIdByAssetId(assetId)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return false
			}
			if userId == user.Id {
				userRbac := entity.UserRbac{
					AssetId: assetId,
					UserId:  user.Id,
					RoleId:  user.RoleId,
				}
				err := service.userRBACRepository.Create(&userRbac, tx)
				if err != nil {
					tx.Rollback()
					return false
				}
			}
		} else {
			userRbac := entity.UserRbac{
				AssetId: assetId,
				UserId:  user.Id,
				RoleId:  user.RoleId,
			}
			err := service.userRBACRepository.Create(&userRbac, tx)
			if err != nil {
				tx.Rollback()
				return false
			}
		}
	}
	tx.Commit()
	return true
}

func (service *AssetsService) UpdateAsset(userId int64, assetId int64, assetName string, purchaseDate time.Time, cost float64, owner *int64, warrantExpiry time.Time, serialNumber string, image *multipart.FileHeader, fileAttachment *multipart.FileHeader, categoryId int64, departmentId int64, status string, maintenance float64) (*entity.Assets, error) {
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
	tx := service.repo.GetDB().Begin()
	asset := entity.Assets{
		Id:                  assetId,
		AssetName:           assetName,
		PurchaseDate:        purchaseDate,
		Cost:                cost,
		Owner:               owner,
		WarrantExpiry:       warrantExpiry,
		SerialNumber:        serialNumber,
		ImageUpload:         &imageUrl,
		FileAttachment:      &fileUrl,
		CategoryId:          categoryId,
		DepartmentId:        departmentId,
		ScheduleMaintenance: &maintenance,
		Status:              status,
	}
	assetUpdated, err := service.repo.UpdateAsset(&asset, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	assetLog := entity.AssetLog{
		Action:        "Update",
		Timestamp:     time.Now(),
		UserId:        userId,
		Asset_id:      assetUpdated.Id,
		ChangeSummary: "Update",
	}
	_, err = service.assertLogRepository.Create(&assetLog, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return assetUpdated, nil
}

func (service *AssetsService) DeleteAsset(userId int64, id int64) error {
	tx := service.repo.GetDB().Begin()
	err := service.repo.DeleteAsset(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	assetLog := entity.AssetLog{
		Action:        "Delete",
		Timestamp:     time.Now(),
		UserId:        userId,
		Asset_id:      id,
		ChangeSummary: "Delete",
	}
	_, err = service.assertLogRepository.Create(&assetLog, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
