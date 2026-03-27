package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PostgreConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type Config struct {
	Postgres *PostgreConfig
	JWT      *JWTConfig
	Redis    *RedisConfig
	HTTPPort string
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func NewPostgreConfig() *PostgreConfig {
	loadEnv()
	return &PostgreConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Username: getEnv("DB_USER", "gamematch_user"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "gamematch_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func NewConfig() *Config {
	loadEnv()
	return &Config{
		Postgres: &PostgreConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "gamematch_user"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "gamematch_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: &JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-me-in-production"),
		},
		Redis: &RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		HTTPPort: getEnv("HTTP_PORT", "8080"),
	}
}
