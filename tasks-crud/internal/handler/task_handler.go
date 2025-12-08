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

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := strings.TrimPrefix(r.URL.Path, "/api/v1")
    
    switch {
    case path == "/tasks" && r.Method == "GET":
        h.getAllTasks(w, r)
    case path == "/tasks" && r.Method == "POST":
        h.createTask(w, r)
    case strings.HasPrefix(path, "/tasks/"):
        h.handleTaskByID(w, r, path)
    default:
        http.NotFound(w, r)
    }
}

func (h *TaskHandler) handleTaskByID(w http.ResponseWriter, r *http.Request, path string) {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) < 2 {
        sendError(w, http.StatusBadRequest, "Invalid path", fmt.Errorf("expected /tasks/{id}"))
        return
    }
    
    id, err := strconv.Atoi(parts[1])
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid task ID", err)
        return
    }
    
    switch r.Method {
    case "GET":
        h.getTaskByID(w, r, id)
    case "PUT":
        h.updateTask(w, r, id)
    case "DELETE":
        h.deleteTask(w, r, id)
    default:
        sendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
    }
}

func (h *TaskHandler) getAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := h.service.GetAllTasks()
    if err != nil {
        sendError(w, http.StatusInternalServerError, "Failed to get tasks", err)
        return
    }
    
    sendJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) getTaskByID(w http.ResponseWriter, r *http.Request, id int) {
    task, err := h.service.GetTaskByID(id)
    if err != nil {
        sendError(w, http.StatusNotFound, "Task not found", err)
        return
    }
    
    sendJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
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

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request, id int) {
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

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
    if err := h.service.DeleteTask(id); err != nil {
        sendError(w, http.StatusNotFound, "Task not found", err)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    if err := json.NewEncoder(w).Encode(data); err != nil {
        fmt.Printf("Failed to encode JSON: %v\n", err)
    }
}

func sendError(w http.ResponseWriter, statusCode int, message string, err error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := map[string]interface{}{
        "error": message,
        "status": statusCode,
    }
    
    if err != nil {
        response["details"] = err.Error()
    }
    
    json.NewEncoder(w).Encode(response)
}