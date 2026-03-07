package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(path string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	return err
}

func Init() {
	if err := InitDB("gesitr.db"); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}
