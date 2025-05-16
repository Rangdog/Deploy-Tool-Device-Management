package config

import (
	"BE_Manage_device/internal/domain/entity"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(DB_DNS), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database. Error:", err)
	}
	err = db.AutoMigrate(&entity.Users{}, &entity.UsersSesions{})
	if err != nil {
		log.Fatal("Error migrate to database. Error:", err)
	}
	return db
}
