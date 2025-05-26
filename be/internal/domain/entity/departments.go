package entity

type Departments struct {
	Id             int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentName string `json:"departmentName"`
	LocationId     int64  `json:"locationId"`

	Location Locations `gorm:"foreignKey:LocationId;references:Id"`
}
