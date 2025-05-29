package entity

import "time"

type AssetLog struct {
	Id                   int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Action               string       `json:"action"`
	Timestamp            time.Time    `json:"timeStamp"`
	UserAssignedId       int64        `json:"userAssignedId"`
	AssetId              int64        `json:"assetId"`
	DepartmentAssignedId *int64       `json:"departmentAssignedId"`
	ChangeSummary        string       `json:"changeSummary"`
	Department           *Departments `gorm:"foreignKey:DepartmentAssignedId;references:Id"`
	Asset                Assets       `gorm:"foreignKey:AssetId;references:Id"`
	User                 Users        `gorm:"foreignKey:UserAssignedId;references:Id"`
}
