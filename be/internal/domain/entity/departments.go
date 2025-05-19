package entity

type Departments struct {
	Id             int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentName string `json:"department_name"`
	LocationId     int64  `json:"location_id"`

	Location Locations `gorm:"foreignKey:LocationId;references:Id"`
}
