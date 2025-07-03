package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	assets "BE_Manage_device/internal/repository/assets"
	bill "BE_Manage_device/internal/repository/bill"
	user "BE_Manage_device/internal/repository/user"
	"BE_Manage_device/pkg/utils"

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

func (service *BillsService) Create(userId int64, assetId int64, description string) (*entity.Bill, error) {
	asset, err := service.assetRepo.GetAssetById(assetId)
	if err != nil {
		return nil, err
	}
	bill := entity.Bill{
		AssetId:     assetId,
		Description: description,
		Amount:      asset.Cost,
		CreateAt:    time.Now(),
		CreateById:  userId,
		CompanyId:   asset.CompanyId,
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
