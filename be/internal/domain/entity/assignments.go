package entity

type Assignments struct {
	Id      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId  *int64 `json:"user_id"`
	AssetId *int64 `json:"asset_id"`
	AssetBy int64  `json:"asset_by"`
}
