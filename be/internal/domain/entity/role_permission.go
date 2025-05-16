package entity

import "time"

type RolePermission struct {
	RoleId       int64      `gorm:"primaryKey" json:"role_id"`
	PermissionId int64      `gorm:"primaryKey" json:"permission_id"`
	AccessLevel  string     `gorm:"type:VARCHAR(50);default:'full'" json:"access_level"`
	Created_at   time.Time  `gorm:"NOT NULL" json:"created_at"`
	Updated_at   *time.Time `json:"updated_at"`
}
