package database

import (
	"log"
	"os"

	gormlog "gorm.io/gorm/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(path string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: gormlog.New(log.Default(), gormlog.Config{
			IgnoreRecordNotFoundError: true,
		}),
	})
	return err
}

func Init() {
	path := os.Getenv("DATABASE_PATH")
	if path == "" {
		path = "gesitr.db"
	}
	if err := InitDB(path); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}
