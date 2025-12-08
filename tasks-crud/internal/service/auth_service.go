package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"tasks-crud/internal/config"
	"tasks-crud/internal/domain"
	"tasks-crud/internal/repository"
)

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

type AuthService struct {
    userRepo repository.UserRepository
    cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        cfg:      cfg,
    }
}

func (s *AuthService) HashPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cfg.BcryptCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(hashedBytes), nil
}

// VerifyPassword - проверка пароля
func (s *AuthService) VerifyPassword(hashedPassword, password string) bool {
    // bcrypt.CompareHashAndPassword сравнивает хэши безопасно (constant-time сравнение)
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

// Register - регистрация нового пользователя
func (s *AuthService) Register(req domain.CreateUserRequest) (*domain.User, error) {
    // Проверяем, существует ли пользователь с таким email
    existingUser, err := s.userRepo.GetByEmail(req.Email)
    if err == nil && existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    // Хэшируем пароль перед сохранением
    hashedPassword, err := s.HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    user := &domain.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        CreatedAt:    time.Now(),
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}

// Login - аутентификация пользователя
func (s *AuthService) Login(req domain.LoginRequest) (*domain.User, error) {
    // Ищем пользователя по email
    user, err := s.userRepo.GetByEmail(req.Email)
    if err != nil {
        // Для безопасности даем одинаковую ошибку, чтобы не раскрывать существование пользователя
        return nil, errors.New("invalid credentials")
    }
    
    // Проверяем пароль (безопасное сравнение)
    if !s.VerifyPassword(user.PasswordHash, req.Password) {
        return nil, errors.New("invalid credentials")
    }
    
    return user, nil
}

// GenerateToken - генерация JWT токена
func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
    // Создаем claims с временем жизни
    claims := Claims{
        UserID:   user.ID,
        Username: user.Username,
        Email:    user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "todo-api",
            Subject:   strconv.Itoa(user.ID),
        },
    }
    
    // Создаем токен с HMAC подписью
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Подписываем токен секретным ключом
    tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %w", err)
    }
    
    return tokenString, nil
}

// ValidateToken - валидация и парсинг JWT токена
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    // Парсим токен с проверкой подписи
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Проверяем алгоритм подписи
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.cfg.JWTSecret), nil
    })
    
    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }
    
    // Извлекаем claims
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token claims")
}