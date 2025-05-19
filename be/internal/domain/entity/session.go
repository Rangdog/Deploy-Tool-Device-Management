package entity

import "time"

type UsersSessions struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64     `gorm:"not null"`
	RefreshToken string    `gorm:"type:text"`
	AccessToken  string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	IsRevoked    bool      `gorm:"default:false" json:"is_revoked"`
}
