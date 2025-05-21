package entity

import (
	"time"
)

type Assets struct {
	Id                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AssetName           string    `json:"assetName"`
	PurchaseDate        time.Time `json:"purchaseDate"`
	Cost                float64   `json:"cost"`
	Owner               *int64    `json:"owner"`
	WarrantExpiry       time.Time `json:"warrantExpiry"`
	Status              string    `json:"status"`
	SerialNumber        string    `json:"serialNumber"`
	FileAttachment      *string   `json:"fileAttachment"`
	ImageUpload         *string   `json:"imageUpload"`
	ScheduleMaintenance *float64  `json:"maintenance"`
	CategoryId          int64     `json:"categoryId"`
	DepartmentId        int64     `json:"departmentId"`

	Category   Categories  `gorm:"foreignKey:CategoryId;references:Id"`
	Department Departments `gorm:"foreignKey:DepartmentId;references:Id"`
	OnwerUser  *Users      `gorm:"foreignKey:Owner;references:Id"`
}
