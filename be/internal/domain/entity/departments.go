package entity

type Departments struct {
	Id             int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentName string `json:"department_name"`
	LocationId     int64  `json:"location_id"`
	UserId         *int64 `json:"-"`

	Location Locations `gorm:"foreignKey:LocationId;references:Id"`
	Head     Users     `gorm:"foreignKey:UserId;references:Id"`
}
