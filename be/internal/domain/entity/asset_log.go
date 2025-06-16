package entity

import "time"

type AssetLog struct {
	Id            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Action        string    `json:"action"`
	Timestamp     time.Time `json:"timeStamp"`
	AssignUserId  *int64    `json:"assignUser"`
	ByUserId      *int64    `json:"byUser"`
	AssetId       int64     `json:"assetId"`
	ChangeSummary string    `json:"changeSummary"`
	ByUser        *Users    `gorm:"foreignKey:ByUserId;references:Id"`
	AssignUser    *Users    `gorm:"foreignKey:AssignUserId;references:Id"`
	Asset         Assets    `gorm:"foreignKey:AssetId;references:Id"`
}
