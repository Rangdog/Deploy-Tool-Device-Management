package entity

type UserRbac struct {
	Id                 int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId             int64 `gorm:"index:unique_userId_AssetId,unique" json:"user_id"`
	RoleId             int64 `json:"role_id"`
	AssetId            int64 `gorm:"index:unique_userId_AssetId,unique" json:"asset_id"`
	NotificationEnable bool  `gorm:"default:true" json:"notification_enable"`

	User Users
}
