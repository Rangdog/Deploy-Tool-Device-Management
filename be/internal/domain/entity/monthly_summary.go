package entity

import "time"

type MonthlySummary struct {
	Id                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Month               int64     `json:"month"`
	Year                int64     `json:"year"`
	TotalAmount         int64     `json:"totalAmount"`
	BillCount           int64     `json:"billCount"`
	AssetCount          int64     `json:"assetCount"`
	TotalCategoryAmount int64     `json:"totalCategoryAmount"`
	GeneratedById       int64     `json:"generatedById"` // user gen id
	GeneratedAt         time.Time `json:"generatedAt"`

	CreateBy Users `gorm:"foreignKey:GeneratedById;references:Id"`
}
