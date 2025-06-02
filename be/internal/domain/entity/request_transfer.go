package entity

type RequestTransfer struct {
	Id          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      int64  `json:"userId"`
	CategoryId  int64  `json:"categoryId"`
	Status      string `json:"status"`
	Description string `json:"description"`

	User     Users      `gorm:"foreignKey:UserId;references:Id"`
	Category Categories `gorm:"foreignKey:CategoryId;references:Id"`
}
