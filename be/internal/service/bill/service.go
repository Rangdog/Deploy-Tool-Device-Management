package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	assets "BE_Manage_device/internal/repository/assets"
	bill "BE_Manage_device/internal/repository/bill"
	user "BE_Manage_device/internal/repository/user"
	"BE_Manage_device/pkg/utils"
	"fmt"
	"mime/multipart"

	"time"
)

type BillsService struct {
	repo      bill.BillsRepository
	assetRepo assets.AssetsRepository
	userRepo  user.UserRepository
}

func NewBillService(repo bill.BillsRepository, assetRepo assets.AssetsRepository, userRepo user.UserRepository) *BillsService {
	return &BillsService{repo: repo, assetRepo: assetRepo, userRepo: userRepo}
}

func (service *BillsService) Create(userId int64, assetId int64, description string, image *multipart.FileHeader, fileAttachment *multipart.FileHeader, status string) (*entity.Bill, error) {
	uploader := utils.NewSupabaseUploader()
	var imageUrl string
	var fileUrl string
	if image != nil {
		imgFile, err := image.Open()
		if err != nil {
			return nil, fmt.Errorf("cannot open image: %w", err)
		}
		defer imgFile.Close()
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), image.Filename)
		imagePath := "bill_image/" + uniqueName
		i, err := uploader.Upload(imagePath, imgFile, image.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		imageUrl = i
	}
	if fileAttachment != nil {
		fileFile, err := fileAttachment.Open()
		if err != nil {
			return nil, fmt.Errorf("cannot open fileAttachment: %w", err)
		}
		defer fileFile.Close()
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileAttachment.Filename)
		filePath := "bill_files/" + uniqueName
		f, err := uploader.Upload(filePath, fileFile, fileAttachment.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		fileUrl = f
	}
	asset, err := service.assetRepo.GetAssetById(assetId)
	if err != nil {
		return nil, err
	}
	bill := entity.Bill{
		AssetId:            assetId,
		Description:        description,
		Amount:             asset.Cost,
		CreateAt:           time.Now(),
		CreateById:         userId,
		CompanyId:          asset.CompanyId,
		FileAttachmentBill: &fileUrl,
		ImageUploadBill:    &imageUrl,
		StatusBill:         status,
	}
	billCreate, err := service.repo.Create(&bill)
	if err != nil {
		return nil, err
	}
	BillGet, err := service.repo.GetByBillNumber(billCreate.BillNumber)
	return BillGet, err
}

func (service *BillsService) GetByBillNumber(billNumber string) (*entity.Bill, error) {
	bill, err := service.repo.GetByBillNumber(billNumber)
	return bill, err
}

func (service *BillsService) Filter(userId int64, BillNumber *string, Status *string, CategoryId *string) ([]dto.BillResponse, error) {
	var filter = filter.BillFilter{
		BillNumber: BillNumber,
		Status:     Status,
		CategoryId: CategoryId,
	}
	users, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	filter.CompanyId = users.CompanyId
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.Bill{}))

	var total int64
	dbFilter.Count(&total)
	var bills []*entity.Bill
	result := dbFilter.Find(&bills)
	if result.Error != nil {
		return nil, result.Error
	}
	billRes := utils.ConvertBillsToResponsesArray(bills)
	return billRes, nil
}

func (service *BillsService) GetAllBillUnpaid(userId int64) ([]*entity.Bill, error) {
	user, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	bills, err := service.repo.GetAllBillUnpaid(user.CompanyId)
	return bills, err
}
func (service *BillsService) UpdatePaid(billNumberStr string) error {
	err := service.repo.UpdatePaid(billNumberStr)
	return err
}
