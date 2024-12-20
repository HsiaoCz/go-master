package database

import (
	"fmt"
	"log"
	"twitter-clone/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(host, user, password, dbname string, port int) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the models
	err = db.AutoMigrate(
		&models.User{},
		&models.Tweet{},
		&models.Follow{},
		&models.Like{},
		&models.Retweet{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
}