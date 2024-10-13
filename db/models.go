package db

import (
	"context"
	"time"

	"gorm.io/gorm"

	"grpc-todo-list/pb"
)

// Task represents a task in the database
type Task struct {
	ID          int32     `gorm:"primaryKey"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	IsCompleted bool      `gorm:"column:is_completed"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

// DB wraps the GORM DB instance
type DB struct {
	*gorm.DB
}

// AddTask adds a new task to the database
func (db *DB) AddTask(ctx context.Context, title, description string) (*pb.TaskResponse, error) {
	task := Task{
		Title:       title,
		Description: description,
		IsCompleted: false,
	}

	if err := db.WithContext(ctx).Create(&task).Error; err != nil {
		return nil, err
	}

	return &pb.TaskResponse{
		Task: &pb.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			IsCompleted: task.IsCompleted,
		},
	}, nil
}

// GetTasks retrieves all tasks from the database
func (db *DB) GetTasks(ctx context.Context) (*pb.GetTasksResponse, error) {
	var tasks []Task

	if err := db.WithContext(ctx).Find(&tasks).Error; err != nil {
		return nil, err
	}

	var pbTasks []*pb.Task
	for _, task := range tasks {
		pbTasks = append(pbTasks, &pb.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			IsCompleted: task.IsCompleted,
		})
	}

	return &pb.GetTasksResponse{Tasks: pbTasks}, nil
}

// CompleteTask marks a task as completed in the database
func (db *DB) CompleteTask(ctx context.Context, id int32) (*pb.TaskResponse, error) {
	var task Task

	if err := db.WithContext(ctx).First(&task, id).Error; err != nil {
		return nil, err
	}

	task.IsCompleted = true
	task.UpdatedAt = time.Now()

	if err := db.WithContext(ctx).Save(&task).Error; err != nil {
		return nil, err
	}

	return &pb.TaskResponse{Task: &pb.Task{Id: task.ID, IsCompleted: task.IsCompleted}}, nil
}
