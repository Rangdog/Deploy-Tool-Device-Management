package entity

import "time"

type AssetLog struct {
	Id            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Action        string    `json:"action"`
	Timestamp     time.Time `json:"timeStamp"`
	UserId        int64     `json:"userId"`
	Asset_id      int64     `json:"assetId"`
	AssignmentId  *int64    `json:"assignmentId"`
	ChangeSummary string    `json:"changeSummary"`
}
