package entity

import (
	"time"
)

type Assets struct {
	Id             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AssetName      string    `json:"assetName"`
	PurchaseDate   time.Time `json:"purchaseDate"`
	Cost           float64   `json:"cost"`
	Owner          *int64    `json:"owner"`
	WarrantExpiry  time.Time `json:"warrantExpiry"`
	Status         string    `gorm:"type:asset_status" json:"status"`
	SerialNumber   string    `json:"serialNumber"`
	FileAttachment *string   `json:"file"`
	ImageUpload    *string   `json:"image"`
	CategoryId     int64     `json:"categoryId"`
	DepartmentId   int64     `json:"departmentId"`
	QrUrl          *string   `json:"qrUrl"`

	Category   Categories  `gorm:"foreignKey:CategoryId;references:Id"`
	Department Departments `gorm:"foreignKey:DepartmentId;references:Id"`
	OnwerUser  *Users      `gorm:"foreignKey:Owner;references:Id"`
}
