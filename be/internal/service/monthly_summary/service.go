package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	bill "BE_Manage_device/internal/repository/bill"
	monthlySummary "BE_Manage_device/internal/repository/monthly_summary"
	user "BE_Manage_device/internal/repository/user"
	"BE_Manage_device/pkg/utils"
)

type MonthlySummaryService struct {
	repo     monthlySummary.MonthlySummaryRepository
	billRepo bill.BillsRepository
	userRepo user.UserRepository
}

func NewMonthlySummaryService(repo monthlySummary.MonthlySummaryRepository, billRepo bill.BillsRepository, userRepo user.UserRepository) *MonthlySummaryService {
	return &MonthlySummaryService{repo: repo, billRepo: billRepo, userRepo: userRepo}
}

func (service *MonthlySummaryService) Filter(userId int64, month, year int64) (*dto.MonthlySummaryResponse, error) {
	var filter = filter.MonthlySummaryFilter{
		Year:  year,
		Month: month,
	}
	users, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	filter.CompanyId = users.CompanyId
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.MonthlySummary{}))
	var MonthlySummary entity.MonthlySummary
	result := dbFilter.First(&MonthlySummary)
	if result.Error != nil {
		return nil, result.Error
	}
	MonthlySummaryRes := utils.ConvertMonthlySummaryToMonthlySummaryRes(MonthlySummary)
	return MonthlySummaryRes, nil
}
