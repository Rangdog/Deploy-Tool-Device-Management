package entity

type RequestTransfer struct {
	Id           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64  `json:"userId"`
	AssetId      int64  `json:"assetId"`
	DepartmentId int64  `json:"departmentId"`
	Status       string `json:"status"`

	User       Users       `gorm:"foreignKey:UserId;references:Id"`
	Asset      Assets      `gorm:"foreignKey:AssetId;references:Id"`
	Department Departments `gorm:"foreignKey:DepartmentId;references:Id"`
}
