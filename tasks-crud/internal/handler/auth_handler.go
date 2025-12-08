package handler

import (
	"encoding/json"
	"net/http"

	"tasks-crud/internal/domain"
	"tasks-crud/internal/service"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

// Register - регистрация нового пользователя
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.CreateUserRequest true "User registration data"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 409 {object} domain.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req domain.CreateUserRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, http.StatusBadRequest, "Invalid JSON", err)
        return
    }
    defer r.Body.Close()
    
    // Базовая валидация
    if len(req.Password) < 8 {
        sendError(w, http.StatusBadRequest, "Password must be at least 8 characters", nil)
        return
    }
    
    // Регистрация пользователя
    user, err := h.authService.Register(req)
    if err != nil {
        status := http.StatusBadRequest
        if err.Error() == "user with this email already exists" {
            status = http.StatusConflict
        }
        sendError(w, status, "Registration failed", err)
        return
    }
    
    // Генерация токена
    token, err := h.authService.GenerateToken(user)
    if err != nil {
        sendError(w, http.StatusInternalServerError, "Failed to generate token", err)
        return
    }
    
    response := domain.AuthResponse{
        Token: token,
        User:  *user,
    }
    
    sendJSON(w, http.StatusCreated, response)
}

// Login - вход пользователя
// @Summary Login user
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "User credentials"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req domain.LoginRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, http.StatusBadRequest, "Invalid JSON", err)
        return
    }
    defer r.Body.Close()
    
    // Аутентификация
    user, err := h.authService.Login(req)
    if err != nil {
        // Для безопасности даем одинаковую ошибку
        sendError(w, http.StatusUnauthorized, "Invalid email or password", nil)
        return
    }
    
    // Генерация токена
    token, err := h.authService.GenerateToken(user)
    if err != nil {
        sendError(w, http.StatusInternalServerError, "Failed to generate token", err)
        return
    }
    
    response := domain.AuthResponse{
        Token: token,
        User:  *user,
    }
    
    sendJSON(w, http.StatusOK, response)
}