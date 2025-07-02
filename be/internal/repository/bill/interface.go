package repository

import (
	"BE_Manage_device/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type BillsRepository interface {
	Create(*entity.Bill) (*entity.Bill, error)
	GetByBillNumber(string) (*entity.Bill, error)
	GetDB() *gorm.DB
	GetAllBillOfMonth(time time.Time, companyId int64) ([]*entity.Bill, error)
}
