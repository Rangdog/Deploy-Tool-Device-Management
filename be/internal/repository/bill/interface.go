package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type BillsRepository interface {
	Create(*entity.Bill) (*entity.Bill, error)
	GetByBillNumber(string) (*entity.Bill, error)
	GetDB() *gorm.DB
}
