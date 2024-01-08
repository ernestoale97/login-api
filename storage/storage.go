package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	db, err := gorm.Open(sqlite.Open("./env/database.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
