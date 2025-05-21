package entity

import "time"

type Roles struct {
	Id          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"type:VARCHAR(250)" json:"title"`
	Slug        string     `gorm:"type:VARCHAR(250);index:unique_slug,unique" json:"slug"`
	Description string     `json:"description"`
	Activated   bool       `gorm:"default:true" json:"activated"`
	Created_at  time.Time  `gorm:"NOT NULL" json:"createdAt"`
	Updated_at  *time.Time `json:"updatedAt"`
	Content     string     `json:"content"`

	Users           []Users          `gorm:"foreignKey:RoleId;references:Id"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleId;references:Id"`
	UserRbacs       []UserRbac       `gorm:"foreignKey:RoleId;references:Id"`
}
