package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"tasks-crud/internal/domain"
	"tasks-crud/internal/service"
)

type TaskHandler struct {
    service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
    return &TaskHandler{
        service: service,
    }
}

// ServeHTTP обрабатывает HTTP запросы (если используете этот метод)
func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := strings.TrimPrefix(r.URL.Path, "/api/v1")
    
    switch {
    case path == "/tasks" && r.Method == "GET":
        h.GetAllTasks(w, r)
    case path == "/tasks" && r.Method == "POST":
        h.CreateTask(w, r)
    case strings.HasPrefix(path, "/tasks/"):
        h.handleTaskByID(w, r, path)
    default:
        http.NotFound(w, r)
    }
}

// GetAllTasks - публичный метод (с большой буквы!)
// @Summary Get all tasks
// @Description Get list of all tasks
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Task
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := h.service.GetAllTasks()
    if err != nil {
        sendError(w, http.StatusInternalServerError, "Failed to get tasks", err)
        return
    }
    
    sendJSON(w, http.StatusOK, tasks)
}

// GetTaskByID - публичный метод
// @Summary Get task by ID
// @Description Get a specific task by ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Security BearerAuth
// @Success 200 {object} domain.Task
// @Failure 404 {object} domain.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
    vars := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(vars) < 2 {
        sendError(w, http.StatusBadRequest, "Invalid path", fmt.Errorf("expected /tasks/{id}"))
        return
    }
    
    id, err := strconv.Atoi(vars[len(vars)-1])
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid task ID", err)
        return
    }
    
    task, err := h.service.GetTaskByID(id)
    if err != nil {
        sendError(w, http.StatusNotFound, "Task not found", err)
        return
    }
    
    sendJSON(w, http.StatusOK, task)
}

// CreateTask - публичный метод
// @Summary Create a new task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body domain.CreateTaskRequest true "Task data"
// @Security BearerAuth
// @Success 201 {object} domain.Task
// @Failure 400 {object} domain.ErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var req domain.CreateTaskRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, http.StatusBadRequest, "Invalid JSON", err)
        return
    }
    defer r.Body.Close()
    
    task, err := h.service.CreateTask(req)
    if err != nil {
        sendError(w, http.StatusBadRequest, "Failed to create task", err)
        return
    }
    
    sendJSON(w, http.StatusCreated, task)
}

// UpdateTask - публичный метод
// @Summary Update a task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body domain.UpdateTaskRequest true "Updated task data"
// @Security BearerAuth
// @Success 200 {object} domain.Task
// @Failure 400 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    vars := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(vars) < 2 {
        sendError(w, http.StatusBadRequest, "Invalid path", fmt.Errorf("expected /tasks/{id}"))
        return
    }
    
    id, err := strconv.Atoi(vars[len(vars)-1])
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid task ID", err)
        return
    }
    
    var req domain.UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, http.StatusBadRequest, "Invalid JSON", err)
        return
    }
    defer r.Body.Close()
    
    task, err := h.service.UpdateTask(id, req)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            sendError(w, http.StatusNotFound, "Task not found", err)
        } else {
            sendError(w, http.StatusBadRequest, "Failed to update task", err)
        }
        return
    }
    
    sendJSON(w, http.StatusOK, task)
}

// DeleteTask - публичный метод
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Security BearerAuth
// @Success 204
// @Failure 404 {object} domain.ErrorResponse
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(vars) < 2 {
        sendError(w, http.StatusBadRequest, "Invalid path", fmt.Errorf("expected /tasks/{id}"))
        return
    }
    
    id, err := strconv.Atoi(vars[len(vars)-1])
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid task ID", err)
        return
    }
    
    if err := h.service.DeleteTask(id); err != nil {
        sendError(w, http.StatusNotFound, "Task not found", err)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

// handleTaskByID - приватный метод для обработки маршрутов с ID
func (h *TaskHandler) handleTaskByID(w http.ResponseWriter, r *http.Request, path string) {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) < 2 {
        sendError(w, http.StatusBadRequest, "Invalid path", fmt.Errorf("expected /tasks/{id}"))
        return
    }
    
    _, err := strconv.Atoi(parts[1])
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid task ID", err)
        return
    }
    
    switch r.Method {
    case "GET":
        h.GetTaskByID(w, r)
    case "PUT":
        h.UpdateTask(w, r)
    case "DELETE":
        h.DeleteTask(w, r)
    default:
        sendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
    }
}

// sendJSON отправляет JSON ответ
func sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    if err := json.NewEncoder(w).Encode(data); err != nil {
        fmt.Printf("Failed to encode JSON: %v\n", err)
    }
}

// sendError отправляет ошибку в JSON формате
func sendError(w http.ResponseWriter, statusCode int, message string, err error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := map[string]interface{}{
        "error":  message,
        "status": statusCode,
    }
    
    if err != nil {
        response["details"] = err.Error()
    }
    
    json.NewEncoder(w).Encode(response)
}