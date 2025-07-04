package entity

import "time"

type Bill struct {
	Id                 int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BillNumber         string    `gorm:"index" json:"billNumber"`
	Amount             float64   `json:"amount"`
	Description        string    `json:"description"`
	CreateAt           time.Time `json:"createAt"`
	CreateById         int64     `json:"createById"`
	AssetId            int64     `gorm:"unique" json:"assetId"`
	StatusBill         string    `json:"statusBill"`
	FileAttachmentBill *string   `json:"fileAttachmentBill"`
	ImageUploadBill    *string   `json:"imageUploadBill"`
	CompanyId          int64     `json:"-"`

	CreateBy Users  `gorm:"foreignKey:CreateById;references:Id"`
	Asset    Assets `gorm:"foreignKey:AssetId;references:Id"`
}
