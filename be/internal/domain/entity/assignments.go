package entity

type Assignments struct {
	Id           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       *int64 `json:"userId"`
	AssetId      int64  `gorm:"index:unique_AssetId,unique" json:"assetId"`
	AssignBy     int64  `json:"assetBy"`
	DepartmentId *int64 `json:"departmentID"`

	UserAssigned Users       `gorm:"foreignKey:UserId;references:Id"`
	UserAssign   Users       `gorm:"foreignKey:AssignBy;references:Id"`
	Asset        Assets      `gorm:"foreignKey:AssetId;references:Id"`
	Department   Departments `gorm:"foreignKey:DepartmentId;references:Id"`
}
