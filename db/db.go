package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database wraps the GORM DB instance
type Database struct {
	*gorm.DB
}

// NewDatabase initializes a new database connection
func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}
