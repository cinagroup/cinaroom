package config

import (
	"os"
	"time"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	CORS      CORSConfig
	CinaToken CinaTokenConfig
}

type CinaTokenConfig struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "multipass"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "multipass-secret-key-change-in-production"),
			ExpireTime: 24 * time.Hour,
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
