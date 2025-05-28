package entity

import "time"

type AssetLog struct {
	Id            int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Action        string       `json:"action"`
	Timestamp     time.Time    `json:"timeStamp"`
	UserId        int64        `json:"userId"`
	AssetId       int64        `json:"assetId"`
	DepartmentId  *int64       `json:"departmentId"`
	ChangeSummary string       `json:"changeSummary"`
	Department    *Departments `gorm:"foreignKey:DepartmentId;references:Id"`
	Asset         Assets       `gorm:"foreignKey:AssetId;references:Id"`
	User          Users        `gorm:"foreignKey:UserId;references:Id"`
}
