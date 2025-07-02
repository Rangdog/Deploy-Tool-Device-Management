package entity

import "time"

type MonthlySummary struct {
	Id                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Month               int64     `json:"month"`
	Year                int64     `json:"year"`
	TotalAmount         float64   `json:"totalAmount"`
	BillCount           int64     `json:"billCount"`
	AssetCount          int64     `json:"assetCount"`
	TotalCategoryAmount string    `json:"totalCategoryAmount"`
	GeneratedAt         time.Time `json:"generatedAt"`
	CompanyId           int64     `json:"-"`
}
