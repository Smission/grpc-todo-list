package service

import (
	"context"
	"fmt"

	"grpc-todo-list/db"
	"grpc-todo-list/pb"
)

type TodoService struct {
	pb.UnimplementedTodoServiceServer
	db *db.Database
}

func NewTodoService(database *db.Database) *TodoService {
	fmt.Println("starting new todo service")
	return &TodoService{db: database}
}

func (s *TodoService) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.TaskResponse, error) {
	task := db.Task{
		Title:       req.Title,
		Description: req.Description,
		IsCompleted: false,
	}

	if err := s.db.Create(&task).Error; err != nil {
		return nil, fmt.Errorf("failed to add task: %w", err)
	}

	return &pb.TaskResponse{Task: &pb.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
	}}, nil
}

// GetTasks retrieves all tasks from the database
func (s *TodoService) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	var tasks []db.Task

	if err := s.db.WithContext(ctx).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
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
func (s *TodoService) CompleteTask(ctx context.Context, req *pb.CompleteTaskRequest) (*pb.TaskResponse, error) {
	var task db.Task

	if err := s.db.WithContext(ctx).First(&task, req.Id).Error; err != nil {
		return nil, fmt.Errorf("failed to find task with ID %d: %w", req.Id, err)
	}

	task.IsCompleted = true

	if err := s.db.WithContext(ctx).Save(&task).Error; err != nil {
		return nil, fmt.Errorf("failed to complete task: %w", err)
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
