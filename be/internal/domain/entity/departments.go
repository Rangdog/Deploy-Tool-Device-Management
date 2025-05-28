package entity

type Departments struct {
	Id             int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentName string `gorm:"uniqueIndex:uniq_dept_location" json:"departmentName"`
	LocationId     int64  `gorm:"uniqueIndex:uniq_dept_location" json:"locationId"`

	Location Locations `gorm:"foreignKey:LocationId;references:Id"`
}
