package entity

import "time"

type Assets struct {
	AssetName      string
	PurchaseDate   time.Time
	Cost           float64
	Owner          int64
	WarrantExpiry  time.Time
	Status         string
	SerialNumber   string
	FileAttachment string
	ImageUpload    string
}
