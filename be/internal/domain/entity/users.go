package entity

type Users struct {
	Id               int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Password         string `gorm:"type:varchar(256);not null" json:"-"`
	FirstName        string `gorm:"type:text;not null;check:(length(first_name)>=2 and length(first_name)<=256)" json:"firstName"`
	LastName         string `gorm:"type:text;not null;check:(length(last_name)>=2 and length(last_name)<=256)" json:"lastName"`
	RoleId           int64  `gorm:"not null" json:"roleId"`
	Email            string `gorm:"unique" json:"email"`
	Token            string `json:"-"`
	IsActive         bool   `json:"isActivate"`
	DepartmentId     *int64 `json:"departmentId"`
	IsHeadDepartment bool   `gorm:"not null;default:false" json:"isHeadDepartment"`
	IsAssetManager   bool   `gorm:"not null;default:false" json:"isAssetManager"`
	CompanyId        int64  `json:"-"`
	CanExport        bool   `gorm:"not null;default:false" json:"canExport"`

	Role       Roles       `gorm:"foreignKey:RoleId;references:Id"`
	Department Departments `gorm:"DepartmentId:RoleId;references:Id"`
}
