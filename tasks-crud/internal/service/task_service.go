package service

import (
	"fmt"
	"strings"

	"tasks-crud/internal/domain"
	"tasks-crud/internal/repository"
)

type TaskService struct {
    repo repository.TaskRepository  
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
    return &TaskService{
        repo: repo,
    }
}

func (s *TaskService) GetAllTasks() ([]domain.Task, error) {
    tasks, err := s.repo.GetAll()
    if err != nil {
        return nil, fmt.Errorf("failed to get tasks: %w", err)
    }
    
    return tasks, nil
}

func (s *TaskService) GetTaskByID(id int) (*domain.Task, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid task id")
    }
    
    task, err := s.repo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("task not found")
    }
    
    return task, nil
}

func (s *TaskService) CreateTask(req domain.CreateTaskRequest) (*domain.Task, error) {
    if err := validateCreateRequest(req); err != nil {
        return nil, err
    }
    
    task := &domain.Task{
        Title:     strings.TrimSpace(req.Title), 
        Completed: false,
    }

    if isDuplicate, err := s.isDuplicateTitle(task.Title); err == nil && isDuplicate {
        return nil, fmt.Errorf("task with title '%s' already exists", task.Title)
    }
    
    if err := s.repo.Create(task); err != nil {
        return nil, fmt.Errorf("failed to create task: %w", err)
    }
    
    return task, nil
}

func (s *TaskService) UpdateTask(id int, req domain.UpdateTaskRequest) (*domain.Task, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid task id: %d", id)
    }
    
    if err := validateUpdateRequest(req); err != nil {
        return nil, err
    }
    
    existingTask, err := s.repo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("task %d not found: %w", id, err)
    }
    
    updatedTask := *existingTask  
    
    if req.Title != nil {
        trimmedTitle := strings.TrimSpace(*req.Title)
        if trimmedTitle == "" {
            return nil, fmt.Errorf("title cannot be empty")
        }
        updatedTask.Title = trimmedTitle
    }
    
    if req.Completed != nil {
        updatedTask.Completed = *req.Completed
    }
    
    if err := s.repo.Update(id, &updatedTask); err != nil {
        return nil, fmt.Errorf("failed to update task: %w", err)
    }
    
    return &updatedTask, nil
}

func (s *TaskService) DeleteTask(id int) error {
    if id <= 0 {
        return fmt.Errorf("invalid task id: %d", id)
    }
    
    if _, err := s.repo.GetByID(id); err != nil {
        return fmt.Errorf("task %d not found: %w", id, err)
    }

    if err := s.repo.Delete(id); err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
    }
    
    fmt.Printf("Task %d deleted\n", id)
    
    return nil
}

func validateCreateRequest(req domain.CreateTaskRequest) error {
    if strings.TrimSpace(req.Title) == "" {
        return fmt.Errorf("title is required")
    }
    
    if len(req.Title) > 200 {
        return fmt.Errorf("title is too long (max 200 characters)")
    }
    
    return nil
}

func validateUpdateRequest(req domain.UpdateTaskRequest) error {
    if req.Title != nil {
        trimmed := strings.TrimSpace(*req.Title)
        if trimmed == "" {
            return fmt.Errorf("title cannot be empty")
        }
        if len(trimmed) > 200 {
            return fmt.Errorf("title is too long (max 200 characters)")
        }
    }
    
    return nil
}

func (s *TaskService) isDuplicateTitle(title string) (bool, error) {
    tasks, err := s.repo.GetAll()
    if err != nil {
        return false, err
    }
    
    cleanTitle := strings.ToLower(strings.TrimSpace(title))
    for _, task := range tasks {
        if strings.ToLower(strings.TrimSpace(task.Title)) == cleanTitle {
            return true, nil
        }
    }
    
    return false, nil
}