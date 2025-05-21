package entity

type Assignments struct {
	Id       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId   *int64 `json:"user_id"`
	AssetId  *int64 `json:"asset_id"`
	AssignBy int64  `json:"asset_by"`

	UserAssigned Users  `gorm:"foreignKey:UserId;references:Id"`
	UserAssign   Users  `gorm:"foreignKey:AssignBy;references:Id"`
	Asset        Assets `gorm:"foreignKey:AssetId;references:Id"`
}
