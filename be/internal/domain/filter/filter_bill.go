package filter

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type BillFilter struct {
	BillNumber *string `form:"billNumber"`
	Status     *string `form:"status"`
	CategoryId *string `form:"categoryId"`
	CompanyId  int64
}

func (f *BillFilter) ApplyFilter(db *gorm.DB) *gorm.DB {
	db.Where("bills.company_id = ?", f.CompanyId)
	if f.Status != nil {
		db = db.Joins("join assets on assets.id = bills.asset_id")
		db = db.Where("status = ?", *f.Status)
	}
	if f.CategoryId != nil {
		db = db.Joins("join assets on assets.id = bills.asset_id")
		db = db.Joins("join categories on categories.id = assets.category_id")
		parsedID, _ := strconv.ParseInt(*f.CategoryId, 10, 64)
		db = db.Where("categories.id = ?", parsedID)
	}
	if f.BillNumber != nil {
		str := "%" + fmt.Sprintf("%v", *f.BillNumber)
		str += "%"
		db = db.Where("LOWER(bills.bill_number) LIKE LOWER(?)", str)
	}
	return db.Preload("CreateBy").Preload("Asset").Preload("CreateBy.Role").Preload("Asset.Category").Preload("Asset.Department").Preload("Asset.Department.Location").Preload("Asset.OnwerUser").Order("bill_number ASC")
}
