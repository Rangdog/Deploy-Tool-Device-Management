package filter

import (
	"gorm.io/gorm"
)

type RequestTransferFilter struct {
	Status *string `form:"status" json:"status"`
	Page   int     `form:"page" json:"page"`
	Limit  int     `form:"limit" json:"limit"`
}

func (f *RequestTransferFilter) ApplyFilter(db *gorm.DB, userId int64) *gorm.DB {
	if f.Status != nil {
		db = db.Where("status = ?", *f.Status)
	}
	return db.Preload("User").Preload("Asset").Preload("Department").Order("request_transfers.id ASC")
}
