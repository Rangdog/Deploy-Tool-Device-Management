package entity

type Users struct {
	Id        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Password  string `gorm:"type:varchar(256);not null" json:"-"`
	FirstName string `gorm:"type:text;not null;check:(length(first_name)>=2 and length(first_name)<=256)" json:"first_name"`
	LastName  string `gorm:"type:text;not null;check:(length(last_name)>=2 and length(last_name)<=256)" json:"last_name"`
	RoleId    int64  `gorm:"not null" json:"role_id"`
	Email     string `gorm:"unique" json:"email"`
	Token     string `json:"-"`
	IsActive  bool   `json:"is_activate"`
	CompanyId int64  `json:"-"`
}
