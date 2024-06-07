package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type User struct {
	ID           uint   `gorm:"primary_key"`
	Username     string `gorm:"unique"`
	Email        string
	PasswordHash string
}

type Image struct {
	ID     uint `gorm:"primary_key"`
	UserID uint
	URL    string
}

func InitDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=testTask port=%s sslmode=disable", host, user, password, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&User{}, &Image{})
	if err != nil {
		return nil
	}
	return db
}
