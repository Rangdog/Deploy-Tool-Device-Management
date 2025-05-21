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

type AssetResponse struct {
	ID             int64              `json:"id"`
	AssetName      string             `json:"assetName"`
	PurchaseDate   string             `json:"purchaseDate"`
	Cost           float64            `json:"cost"`
	Owner          OwnerResponse      `json:"owner,omitempty"`
	WarrantExpiry  string             `json:"warrantExpiry"`
	Status         string             `json:"status"`
	SerialNumber   string             `json:"serialNumber"`
	FileAttachment string             `json:"fileAttachment"`
	ImageUpload    string             `json:"imageUpload"`
	Maintenance    string             `json:"maintenance,omitempty"`
	Category       CategoryResponse   `json:"category"`
	Department     DepartmentResponse `json:"department"`
}

type CategoryResponse struct {
	ID           int64  `json:"id"`
	CategoryName string `json:"categoryName"`
}

type DepartmentResponse struct {
	ID             int64            `json:"id"`
	DepartmentName string           `json:"departmentName"`
	Location       LocationResponse `json:"location"`
}

type LocationResponse struct {
	ID           int64  `json:"id"`
	LocationName string `json:"locationAddress"`
}

type OwnerResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
