package entity

import "time"

type UsersSessions struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64     `gorm:"not null;index:idx_userid_revoked"`
	RefreshToken string    `gorm:"type:text"`
	AccessToken  string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"not null" json:"createdAt"`
	ExpiresAt    time.Time `gorm:"not null" json:"expiresAt"`
	IsRevoked    bool      `gorm:"default:false;index:idx_userid_revoked" json:"isRevoked"`
}
