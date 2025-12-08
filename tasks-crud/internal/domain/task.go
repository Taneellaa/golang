package domain

import "time"

// Task модель задачи
// @Description Модель задачи пользователя
type Task struct {
    ID        int       `json:"id" example:"1"`
    Title     string    `json:"title" example:"Learn Go"`
    Completed bool      `json:"completed" example:"false"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// CreateTaskRequest запрос на создание задачи
// @Description Запрос для создания новой задачи
type CreateTaskRequest struct {
    Title string `json:"title" binding:"required" example:"Learn Go programming"`
}

// UpdateTaskRequest запрос на обновление задачи
// @Description Запрос для обновления существующей задачи
type UpdateTaskRequest struct {
    Title     *string `json:"title,omitempty" example:"Learn Go programming"`    
    Completed *bool   `json:"completed,omitempty" example:"true"`  
}