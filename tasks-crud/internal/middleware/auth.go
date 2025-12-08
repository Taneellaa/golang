package middleware

import (
	"context"
	"net/http"
	"strings"

	"tasks-crud/internal/service"
)

type contextKey string

const (
    UserIDKey   contextKey = "user_id"
    UsernameKey contextKey = "username"
)

// JWTAuthMiddleware - middleware для проверки JWT токена
func JWTAuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Получаем заголовок Authorization
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, `{"error": "Authorization header is required"}`, http.StatusUnauthorized)
                return
            }
            
            // Проверяем формат "Bearer <token>"
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, `{"error": "Invalid authorization format. Use: Bearer <token>"}`, http.StatusUnauthorized)
                return
            }
            
            tokenString := parts[1]
            
            // Валидируем токен
            claims, err := authService.ValidateToken(tokenString)
            if err != nil {
                http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
                return
            }
            
            // Добавляем данные пользователя в контекст
            ctx := r.Context()
            ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
            ctx = context.WithValue(ctx, UsernameKey, claims.Username)
            
            // Передаем запрос дальше с обновленным контекстом
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(UserIDKey).(int)
    return userID, ok
}

func GetUsernameFromContext(ctx context.Context) (string, bool) {
    username, ok := ctx.Value(UsernameKey).(string)
    return username, ok
}