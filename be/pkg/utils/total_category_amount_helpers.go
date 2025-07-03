package utils

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	bill "BE_Manage_device/internal/repository/bill"
	company "BE_Manage_device/internal/repository/company"
	monthlySummary "BE_Manage_device/internal/repository/monthly_summary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func ConvertTCAToStr(TCA map[string]float64) string {
	var str string
	for k, v := range TCA {
		vStr := fmt.Sprintf("%f", v)
		str += k + ":" + vStr + " "
	}
	return str
}

func ConvertStrToTCA(TCAStr string) []*dto.TotalCategoryAmountResponse {
	var TCARes []*dto.TotalCategoryAmountResponse
	TCA := strings.Split(TCAStr, ` `)
	for _, tca := range TCA {
		tcaSplit := strings.Split(tca, `:`)
		amountParse, _ := strconv.ParseFloat(tcaSplit[1], 64)
		tcaRes := dto.TotalCategoryAmountResponse{
			CategoryName: tcaSplit[0],
			Amount:       amountParse,
		}
		TCARes = append(TCARes, &tcaRes)
	}
	return TCARes
}

func ConvertMonthlySummaryToMonthlySummaryRes(MonthlySummary entity.MonthlySummary) *dto.MonthlySummaryResponse {
	var TotalCategoryAmountRes []*dto.TotalCategoryAmountResponse
	if MonthlySummary.TotalCategoryAmount != "" {
		TotalCategoryAmountRes = ConvertStrToTCA(MonthlySummary.TotalCategoryAmount)
	} else {
		TotalCategoryAmountRes = []*dto.TotalCategoryAmountResponse{}
	}
	return &dto.MonthlySummaryResponse{
		Month:               MonthlySummary.Month,
		Year:                MonthlySummary.Year,
		TotalAmount:         MonthlySummary.TotalAmount,
		BillCount:           MonthlySummary.BillCount,
		AssetCount:          MonthlySummary.AssetCount,
		TotalCategoryAmount: TotalCategoryAmountRes,
		GeneratedAt:         MonthlySummary.GeneratedAt,
	}
}

func CreateSummary(MonthlySummaryRepo monthlySummary.MonthlySummaryRepository, billRepo bill.BillsRepository, companyRepo company.CompanyRepository) {
	company, err := companyRepo.GetAllCompany()
	if err != nil {
		logrus.Infof("Happen error when create summary at: %v", time.Now())
		return
	}
	time := time.Now()
	month := time.Month()
	year := time.Year()
	for _, c := range company {
		var totalAmount float64
		var BillCount int64
		var AssetCount int64
		TotalCategoryAmount := make(map[string]float64)
		bills, err := billRepo.GetAllBillOfMonth(time, c.Id)
		if err != nil {
			logrus.Infof("Happen error get bill of company: %v", c.CompanyName)
			continue
		}
		for _, b := range bills {
			BillCount += 1
			AssetCount += 1
			val, ok := TotalCategoryAmount[b.Asset.Category.CategoryName]
			if ok {
				TotalCategoryAmount[b.Asset.Category.CategoryName] = (val + b.Asset.Cost)
			} else {
				TotalCategoryAmount[b.Asset.Category.CategoryName] = b.Asset.Cost
			}
			totalAmount += b.Asset.Cost
		}
		TotalCategoryAmountStr := ConvertTCAToStr(TotalCategoryAmount)

		monthlySummary := entity.MonthlySummary{
			Month:               int64(month),
			Year:                int64(year),
			GeneratedAt:         time,
			TotalCategoryAmount: TotalCategoryAmountStr,
			AssetCount:          AssetCount,
			BillCount:           BillCount,
			TotalAmount:         totalAmount,
			CompanyId:           c.Id,
		}
		_, err = MonthlySummaryRepo.Create(&monthlySummary)
		if err != nil {
			logrus.Infof("Happen error when create summary for company: %v", c.CompanyName)
			continue
		}
	}

}
