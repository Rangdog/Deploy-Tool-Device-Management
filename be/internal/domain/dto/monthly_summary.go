package dto

import (
	"time"
)

type MonthlySummaryResponse struct {
	Month               int64                          `json:"month"`
	Year                int64                          `json:"year"`
	TotalAmount         float64                        `json:"totalAmount"`
	BillCount           int64                          `json:"billCount"`
	AssetCount          int64                          `json:"assetCount"`
	TotalCategoryAmount []*TotalCategoryAmountResponse `json:"totalCategoryAmount"`
	GeneratedAt         time.Time                      `json:"generatedAt"`
}

type TotalCategoryAmountResponse struct {
	CategoryName string  `json:"categoryName"`
	Amount       float64 `json:"amount"`
}
