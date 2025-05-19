package entity

type Categories struct {
	Id           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryName string `json:"category_name"`
}
