package dto

import (
	"time"
)

type BillCreateRequest struct {
	AssetId     int64  `json:"assetId"`
	Description string `json:"description"`
}

type BillResponse struct {
	BillNumber  string        `json:"billNumber"`
	Amount      float64       `json:"amount"`
	Description string        `json:"description"`
	CreateAt    time.Time     `json:"createAt"`
	Asset       AssetResponse `json:"assets"`
	CreateBy    UserResponse  `json:"createBy"`
}
