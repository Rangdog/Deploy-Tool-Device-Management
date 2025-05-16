package entity

import "time"

type Permission struct {
	Id          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"type:VARCHAR(250)" json:"title"`
	Slug        string     `gorm:"type:VARCHAR(250);uniqueIndex" json:"slug"`
	Description string     `json:"description"`
	Activated   bool       `gorm:"default:true" json:"activated"`
	Created_at  time.Time  `grom:"NOT NULL" json:"created_at"`
	Updated_at  *time.Time `json:"updated_at"`
	Content     string     `json:"content"`

	RolePermissions []RolePermission `gorm:"foreignKey:PermissionId;references:Id"`
}
