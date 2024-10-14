package service_test

import (
	"context"
	"testing"

	// Adjust import path as needed

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"grpc-todo-list/pb"      // Adjust import path as needed
	"grpc-todo-list/service" // Adjust import path as needed
)

func TestGetTasks(t *testing.T) {
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

	// Create a new instance of TodoService with the mocked database
	todoService := &service.TodoService{db: gormDB}

	// Prepare the expected tasks
	expectedTasks := []db.Task{
		{ID: 1, Title: "Task 1", Description: "Description 1", IsCompleted: false},
		{ID: 2, Title: "Task 2", Description: "Description 2", IsCompleted: true},
	}

	// Set up the mock to expect a call to Find
	rows := sqlmock.NewRows([]string{"id", "title", "description", "is_completed"}).
		AddRow(expectedTasks[0].ID, expectedTasks[0].Title, expectedTasks[0].Description, expectedTasks[0].IsCompleted).
		AddRow(expectedTasks[1].ID, expectedTasks[1].Title, expectedTasks[1].Description, expectedTasks[1].IsCompleted)

	mock.ExpectQuery("SELECT (.+) FROM `tasks`").WillReturnRows(rows)

	// Create a request
	req := &pb.GetTasksRequest{}

	// Call GetTasks
	resp, err := todoService.GetTasks(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(expectedTasks), len(resp.Tasks))

	for i, task := range resp.Tasks {
		assert.Equal(t, expectedTasks[i].ID, task.Id)
		assert.Equal(t, expectedTasks[i].Title, task.Title)
		assert.Equal(t, expectedTasks[i].Description, task.Description)
		assert.Equal(t, expectedTasks[i].IsCompleted, task.IsCompleted)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
