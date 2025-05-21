package entity

type Locations struct {
	Id           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	LocationName string `json:"locationAddress"`
}
