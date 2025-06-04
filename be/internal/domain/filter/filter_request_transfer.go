package filter

import (
	"gorm.io/gorm"
)

type RequestTransferFilter struct {
	Status *string `form:"status" json:"status"`
}

func (f *RequestTransferFilter) ApplyFilter(db *gorm.DB, userId int64) *gorm.DB {
	if f.Status != nil && *f.Status != "" {
		db = db.Where("status = ?", *f.Status)
	}
	return db.Preload("User").Preload("User.Department").Preload("Category").Order("request_transfers.id ASC")
}
