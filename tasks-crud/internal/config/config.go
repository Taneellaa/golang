package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
    Port         int
    Env          string
    JWTSecret    string
    JWTExpiry    time.Duration
    BcryptCost   int 
}

func Load() *Config {
    port := getEnvAsInt("PORT", 8080)
    env := getEnv("ENV", "development")
    jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
    jwtExpiry := getEnvAsDuration("JWT_EXPIRY", 24*time.Hour) 
    bcryptCost := getEnvAsInt("BCRYPT_COST", 12) 
    
    if env == "production" && jwtSecret == "your-secret-key-change-in-production" {
        panic("JWT_SECRET must be set in production environment")
    }
    
    return &Config{
        Port:       port,
        Env:        env,
        JWTSecret:  jwtSecret,
        JWTExpiry:  jwtExpiry,
        BcryptCost: bcryptCost,
    }
}

// Get возвращает глобальную конфигурацию
func Get() *Config {
    return cfg
}

var cfg = Load()


func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value, exists := os.LookupEnv(key); exists {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}