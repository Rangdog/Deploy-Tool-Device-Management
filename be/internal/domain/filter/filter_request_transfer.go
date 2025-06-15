package filter

import (
	"gorm.io/gorm"
)

type RequestTransferFilter struct {
	Status *string `form:"status" json:"status"`
	DepId *int64 
}

func (f *RequestTransferFilter) ApplyFilter(db *gorm.DB, userId int64) *gorm.DB {
	if f.DepId != nil{
		db.Joins("JOIN users ON users.id = request_transfers.user_id").Where("users.department_id != ?", *f.DepId)
	}
	if f.Status != nil && *f.Status != "" {
		db = db.Where("status = ?", *f.Status)
	}
	return db.Preload("User").Preload("User.Department").Preload("Category").Order("request_transfers.id ASC")
}
