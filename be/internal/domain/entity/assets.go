package entity

import (
	"time"
)

type Assets struct {
	Id                  int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	AssetName           string     `json:"asset_name"`
	PurchaseDate        time.Time  `json:"purchase_date"`
	Cost                float64    `json:"cost"`
	Owner               *int64     `json:"Owner"`
	WarrantExpiry       time.Time  `json:"warrant_expiry"`
	Status              string     `json:"status"`
	SerialNumber        string     `json:"serial_number"`
	FileAttachment      *string    `json:"file_attachment"`
	ImageUpload         *string    `json:"image_upload"`
	ScheduleMaintenance *time.Time `json:"maintenance"`
	CategoryId          int64      `json:"category_id"`
	DepartmentId        *int64     `json:"department_id"`

	Category   Categories  `gorm:"foreignKey:CategoryId;references:Id"`
	Department Departments `gorm:"foreignKey:DepartmentId;references:Id"`
}
