package entity

import "time"

type MaintenanceNotifications struct {
	Id         int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ScheduleId int64
	NotifyDate time.Time

	MaintenanceSchedule MaintenanceSchedules `gorm:"foreignKey:ScheduleId;references:Id"`
}
