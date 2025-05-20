package dto

import "time"

type AssetCreateRequest struct {
	AssetName     string    `json:"asset_name" binding:"required"`
	PurchaseDate  time.Time `json:"purchase_date"  binding:"required"`
	Cost          float64   `json:"cost"  binding:"required"`
	Owner         *int64    `json:"owner"`
	WarrantExpiry time.Time `json:"warrant_expiry"  binding:"required"`
	Status        string    `json:"status" binding:"required"`
	SerialNumber  string    `json:"serial_number"  binding:"required"`
	CategoryId    int64     `json:"category_id"  binding:"required"`
	DepartmentId  *int64    `json:"department_id"`
}
