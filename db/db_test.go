package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestNewDatabase tests the NewDatabase function
func TestNewDatabase(t *testing.T) {
	// Create a new sqlmock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	// Create a GORM DB instance using the mocked DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	// Set expectations for the database interactions
	mock.ExpectPing()
	// Mock the call to SELECT VERSION()
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.0"))

	// Initialize your Database struct
	database := &Database{gormDB}

	// Test the database connection by executing a simple query
	if err := database.DB.Exec("SELECT 1").Error; err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
