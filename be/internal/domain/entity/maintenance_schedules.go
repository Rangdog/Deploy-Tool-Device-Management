package entity

import "time"

type MaintenanceSchedules struct {
	Id        int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	AssetId   int64
	StartDate time.Time
	EndDate   time.Time

	Asset Assets `gorm:"foreignKey:AssetId;references:Id"`
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}
