package domain

import "time"

type Task struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
    Title string `json:"title" binding:"required"`
}

type UpdateTaskRequest struct {
    Title     *string `json:"title,omitempty"`    
    Completed *bool   `json:"completed,omitempty"`  
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}