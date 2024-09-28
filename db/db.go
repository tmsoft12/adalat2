package config

import (
	"log"
	model "tm/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	database.AutoMigrate(&model.BannerSchema{}, &model.EmployerSchema{}, &model.MediaSchema{}, &model.NewsSchema{})
	DB = database
}
