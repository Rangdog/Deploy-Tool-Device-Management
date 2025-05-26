package entity

import "time"

type Notifications struct {
	Id         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Content    *string    `json:"content"`
	Status     *string    `json:"status"`
	UserId     *int64     `json:"userId"`
	Type       *string    `json:"type"`
	AssetId    *int64     `json:"assetId"`
	NotifyDate *time.Time `json:"notifyDate"`
	User       Users      `gorm:"foreignKey:UserId;references:Id"`
	Asset      Assets     `gorm:"foreignKey:AssetId;references:Id"`
}
