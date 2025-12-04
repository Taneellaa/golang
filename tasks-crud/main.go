package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"tasks-crud/models"
	"time"
)

var (
    tasks     = make(map[int]models.Task)  
    currentID = 1                   
    mu        sync.RWMutex          
)

func main() {
    tasks[1] = models.Task{
        ID:        1,
        Title:     "Выучить основы Go",
        Completed: false,
        CreatedAt: time.Now(),
    }
    tasks[2] = models.Task{
        ID:        2,
        Title:     "Написать первое API",
        Completed: true,
        CreatedAt: time.Now(),
    }
    currentID = 3

    http.HandleFunc("/tasks", handleTasks)      
    http.HandleFunc("/tasks/", handleTaskById)  

    fmt.Println("🚀 Сервер запущен: http://localhost:8080")
    fmt.Println("📌 Эндпоинты:")
    fmt.Println("   GET    /tasks      - список всех задач")
    fmt.Println("   POST   /tasks      - создать задачу")
    fmt.Println("   GET    /tasks/{id} - получить задачу")
    fmt.Println("   PUT    /tasks/{id} - обновить задачу")
    fmt.Println("   DELETE /tasks/{id} - удалить задачу")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    switch r.Method {
    case "GET":
        getAllTasks(w)
    case "POST":
        createTask(w, r)
    default:
        errorResponse(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func handleTaskById(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    idStr := r.URL.Path[len("/tasks/"):]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        errorResponse(w, "Неверный ID задачи", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case "GET":
        getTaskById(w, id)
    case "PUT":
        updateTask(w, r, id)
    case "DELETE":
        deleteTask(w, id)
    default:
        errorResponse(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getAllTasks(w http.ResponseWriter) {
    mu.RLock()
    defer mu.RUnlock()

    taskList := make([]models.Task, 0, len(tasks))
    for _, task := range tasks {
        taskList = append(taskList, task)
    }

    json.NewEncoder(w).Encode(taskList)
}

func createTask(w http.ResponseWriter, r *http.Request) {
    var task models.Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        errorResponse(w, "Неверный JSON", http.StatusBadRequest)
        return
    }

    if task.Title == "" {
        errorResponse(w, "Поле 'title' обязательно", http.StatusBadRequest)
        return
    }

    mu.Lock() 
    defer mu.Unlock()

    task.ID = currentID
    task.CreatedAt = time.Now()
    tasks[currentID] = task
    currentID++

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

func getTaskById(w http.ResponseWriter, id int) {
    mu.RLock()
    defer mu.RUnlock()

    task, exists := tasks[id]
    if !exists {
        errorResponse(w, "Задача не найдена", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request, id int) {
    mu.Lock()
    defer mu.Unlock()

    task, exists := tasks[id]
    if !exists {
        errorResponse(w, "Задача не найдена", http.StatusNotFound)
        return
    }

    var updates struct {
        Title     *string `json:"title"`     
        Completed *bool   `json:"completed"`
    }

    err := json.NewDecoder(r.Body).Decode(&updates)
    if err != nil {
        errorResponse(w, "Неверный JSON", http.StatusBadRequest)
        return
    }

    if updates.Title != nil {
        task.Title = *updates.Title
    }
    if updates.Completed != nil {
        task.Completed = *updates.Completed
    }

    tasks[id] = task
    json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, id int) {
    mu.Lock()
    defer mu.Unlock()

    _, exists := tasks[id]
    if !exists {
        errorResponse(w, "Задача не найдена", http.StatusNotFound)
        return
    }

    delete(tasks, id)
    w.WriteHeader(http.StatusNoContent) 
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]string{
        "error":   message,
        "status":  http.StatusText(statusCode),
        "code":    strconv.Itoa(statusCode),
    })
}