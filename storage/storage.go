package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"login_api/pkg/config"
)

var db *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	settings, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(settings.AppDatabase), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
