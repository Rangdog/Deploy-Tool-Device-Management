package entity

type Assignments struct {
	Id       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId   *int64 `json:"userId"`
	AssetId  *int64 `json:"assetId"`
	AssignBy int64  `json:"assetBy"`

	UserAssigned Users  `gorm:"foreignKey:UserId;references:Id"`
	UserAssign   Users  `gorm:"foreignKey:AssignBy;references:Id"`
	Asset        Assets `gorm:"foreignKey:AssetId;references:Id"`
}
