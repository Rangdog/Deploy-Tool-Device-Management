package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type MonthlySummaryRepository interface {
	Create(*entity.MonthlySummary) (*entity.MonthlySummary, error)
	GetDB() *gorm.DB
}
