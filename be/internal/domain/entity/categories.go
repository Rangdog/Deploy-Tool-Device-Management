package entity

type Categories struct {
	Id           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryName string `gorm:"unique" json:"categoryName"`
}
