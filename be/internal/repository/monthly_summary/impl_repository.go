package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type PostgreSQLMonthlySummary struct {
	db *gorm.DB
}

func NewPostgreSQLMonthlySummary(db *gorm.DB) MonthlySummaryRepository {
	return &PostgreSQLMonthlySummary{db: db}
}

func (r *PostgreSQLMonthlySummary) Create(monthlySummary *entity.MonthlySummary) (*entity.MonthlySummary, error) {
	result := r.db.Create(monthlySummary)
	return monthlySummary, result.Error
}

func (r *PostgreSQLMonthlySummary) GetDB() *gorm.DB {
	return r.db
}
