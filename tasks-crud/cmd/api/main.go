package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "tasks-crud/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"tasks-crud/internal/config"
	"tasks-crud/internal/domain"
	"tasks-crud/internal/repository"
	"tasks-crud/internal/service"
)

type TaskHandler struct {
    service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
    return &TaskHandler{service: service}
}


func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := h.service.GetAllTasks()
    if err != nil {
        sendError(w, http.StatusInternalServerError, "Failed to get tasks", err)
        return
    }
    sendJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
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

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
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
        sendError(w, http.StatusBadRequest, "Failed to update task", err)
        return
    }
    
    sendJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
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

func HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status":    "ok",
        "timestamp": time.Now().Format(time.RFC3339),
        "service":   "todo-api",
    })
}

func sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, statusCode int, message string, err error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(domain.ErrorResponse{
        Error:   message,
        Details: err.Error(),
    })
}

func main() {
    fmt.Println("🚀 Запуск Todo API со Swagger...")
    
    cfg := config.Load()
    fmt.Printf("📋 Конфигурация:\n   Порт: %d\n", cfg.Port)
    
    taskRepo := repository.NewInMemoryTaskRepository()
    taskService := service.NewTaskService(taskRepo)
    taskHandler := NewTaskHandler(taskService)
    
    router := mux.NewRouter()
    
    api := router.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
    api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
    api.HandleFunc("/tasks/{id}", taskHandler.GetTaskByID).Methods("GET")
    api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
    api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")
    
    router.HandleFunc("/health", HealthCheck).Methods("GET")
    
    router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
        httpSwagger.URL("/swagger/doc.json"),
        httpSwagger.DeepLinking(true),
        httpSwagger.DocExpansion("list"),
        httpSwagger.DomID("swagger-ui"),
    ))
    
    router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
    })
    
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/swagger/index.html", http.StatusSeeOther)
    })
    
    addr := fmt.Sprintf(":%d", cfg.Port)
    server := &http.Server{
        Addr:         addr,
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    fmt.Printf("🌐 Сервер запущен на http://localhost:%d\n", cfg.Port)
    fmt.Printf("📚 Swagger UI: http://localhost:%d/swagger/index.html\n", cfg.Port)
    fmt.Printf("📖 API документация: http://localhost:%d/docs\n", cfg.Port)
    fmt.Println("🛑 Для остановки нажмите Ctrl+C")
    
    log.Fatal(server.ListenAndServe())
}