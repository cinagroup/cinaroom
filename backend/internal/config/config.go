package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	CORS      CORSConfig
	CinaToken CinaTokenConfig
	Log       LogConfig
}

// CinaTokenConfig holds CinaToken OAuth settings.
type CinaTokenConfig struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       string
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         string
	Mode         string // debug | release
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	Schema          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

// CORSConfig holds CORS middleware settings.
type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string // debug | info | warn | error
	Format string // json | text
}

// Load reads configuration from environment variables with sensible defaults.
// It automatically loads a .env file if present (via godotenv or manual parsing).
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Mode:         getEnv("SERVER_MODE", "debug"),
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "cinaroom"),
			Password:        getEnv("DB_PASSWORD", "cinaroom"),
			DBName:          getEnv("DB_NAME", "cinatoken"),
			Schema:          getEnv("DB_SCHEMA", "cinaroom"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME", 3600)) * time.Second,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "cinaroom-secret-key-change-in-production"),
			ExpireTime: time.Duration(getEnvInt("JWT_EXPIRE_HOURS", 24)) * time.Hour,
		},
		CORS: CORSConfig{
			AllowOrigins: getEnvSlice("CORS_ALLOW_ORIGINS", []string{"*"}),
			AllowMethods: getEnvSlice("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowHeaders: getEnvSlice("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CinaToken-Token"}),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		CinaToken: CinaTokenConfig{
			BaseURL:      getEnv("CINATOKEN_BASE_URL", "https://cinatoken.com"),
			ClientID:     getEnv("CINATOKEN_CLIENT_ID", ""),
			ClientSecret: getEnv("CINATOKEN_CLIENT_SECRET", ""),
			RedirectURI:  getEnv("CINATOKEN_REDIRECT_URI", "http://localhost:3000/oauth/callback"),
			Scopes:       getEnv("CINATOKEN_SCOPES", "user:read user:email"),
		},
	}
}

// getEnv returns the environment variable value or the default.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns the environment variable as an integer or the default.
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
		slog.Warn("invalid integer env var, using default", "key", key, "value", value, "default", defaultValue)
	}
	return defaultValue
}

// getEnvSlice returns a comma-separated env var as a string slice, or the default.
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		if value == "*" {
			return []string{"*"}
		}
		return splitCSV(value)
	}
	return defaultValue
}

// splitCSV splits a comma-separated string into trimmed slices.
func splitCSV(s string) []string {
	var result []string
	for _, v := range splitString(s, ",") {
		result = append(result, v)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			part := trimSpaces(s[start:i])
			if part != "" {
				result = append(result, part)
			}
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	if start < len(s) {
		part := trimSpaces(s[start:])
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func trimSpaces(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
