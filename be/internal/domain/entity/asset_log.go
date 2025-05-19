package entity

import "time"

type AssetLog struct {
	Id            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Action        string    `json:"action"`
	Timestamp     time.Time `json:"time_stamp"`
	UserId        int64     `json:"user_id"`
	Asset_id      int64     `json:"asset_id"`
	AssignmentId  *int64    `json:"assignment_id"`
	ChangeSummary string    `json:"change_summary"`
}
