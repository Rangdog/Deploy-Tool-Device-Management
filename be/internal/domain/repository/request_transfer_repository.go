package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type RequestTransferRepository interface {
	Create(*entity.RequestTransfer) (*entity.RequestTransfer, error)
	UpdateStatusConfirm(id int64, tx *gorm.DB) (*entity.RequestTransfer, error)
	UpdateStatusDeny(id int64, tx *gorm.DB) (*entity.RequestTransfer, error)
	GetDB() *gorm.DB
	GetRequestTransferById(id int64) (*entity.RequestTransfer, error)
}
