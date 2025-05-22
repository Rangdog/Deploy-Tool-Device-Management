package entity

type Notifications struct {
	Id      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Content string `json:"content"`
	Status  string `json:"status"`
	UserId  int64  `json:"userId"`

	User Users `gorm:"foreignKey:UserId;references:Id"`
}
