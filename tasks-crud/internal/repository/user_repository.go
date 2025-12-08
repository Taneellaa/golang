package repository

import (
	"fmt"
	"sync"
	"time"

	"tasks-crud/internal/domain"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
    Create(user *domain.User) error
    GetByEmail(email string) (*domain.User, error)
    GetByID(id int) (*domain.User, error)
    GetAll() ([]domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
}

// InMemoryUserRepository реализация in-memory репозитория пользователей
type InMemoryUserRepository struct {
    users     map[int]domain.User
    emails    map[string]int // для быстрого поиска по email
    currentID int
    mu        sync.RWMutex
}

// NewInMemoryUserRepository создает новый репозиторий пользователей
func NewInMemoryUserRepository() *InMemoryUserRepository {
    repo := &InMemoryUserRepository{
        users:     make(map[int]domain.User),
        emails:    make(map[string]int),
        currentID: 1,
    }
    
    // Создаем тестового пользователя для демонстрации
    testUser := domain.User{
        ID:           1,
        Username:     "admin",
        Email:        "admin@example.com",
        PasswordHash: "$2a$12$N9qo8uLOickgx2ZMRZoMyeMRZDzX4fLmHwQc6UimjX8cK7s2p6qVy", // password: admin123
        CreatedAt:    time.Now(),
    }
    
    repo.users[1] = testUser
    repo.emails["admin@example.com"] = 1
    repo.currentID = 2
    
    return repo
}

// Create создает нового пользователя
func (r *InMemoryUserRepository) Create(user *domain.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Проверяем уникальность email
    if _, exists := r.emails[user.Email]; exists {
        return fmt.Errorf("user with email %s already exists", user.Email)
    }

    // Проверяем уникальность username (опционально)
    for _, existingUser := range r.users {
        if existingUser.Username == user.Username {
            return fmt.Errorf("user with username %s already exists", user.Username)
        }
    }

    user.ID = r.currentID
    user.CreatedAt = time.Now()
    r.users[r.currentID] = *user
    r.emails[user.Email] = r.currentID
    r.currentID++

    return nil
}

// GetByEmail ищет пользователя по email
func (r *InMemoryUserRepository) GetByEmail(email string) (*domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    userID, exists := r.emails[email]
    if !exists {
        return nil, fmt.Errorf("user with email %s not found", email)
    }

    user, exists := r.users[userID]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }

    return &user, nil
}

// GetByID ищет пользователя по ID
func (r *InMemoryUserRepository) GetByID(id int) (*domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    user, exists := r.users[id]
    if !exists {
        return nil, fmt.Errorf("user with id %d not found", id)
    }

    return &user, nil
}

// GetAll возвращает всех пользователей
func (r *InMemoryUserRepository) GetAll() ([]domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    users := make([]domain.User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }

    return users, nil
}

// Update обновляет данные пользователя
func (r *InMemoryUserRepository) Update(user *domain.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.users[user.ID]; !exists {
        return fmt.Errorf("user with id %d not found", user.ID)
    }

    // Если меняется email, проверяем уникальность нового email
    oldUser := r.users[user.ID]
    if oldUser.Email != user.Email {
        if _, exists := r.emails[user.Email]; exists {
            return fmt.Errorf("user with email %s already exists", user.Email)
        }
        delete(r.emails, oldUser.Email)
        r.emails[user.Email] = user.ID
    }

    r.users[user.ID] = *user
    return nil
}

// Delete удаляет пользователя
func (r *InMemoryUserRepository) Delete(id int) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    user, exists := r.users[id]
    if !exists {
        return fmt.Errorf("user with id %d not found", id)
    }

    delete(r.emails, user.Email)
    delete(r.users, id)

    return nil
}