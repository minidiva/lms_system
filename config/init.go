package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env не найден, читаю переменные окружения напрямую")
	}

	cfg := &Config{}

	// Database configuration
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.Name = getEnv("DB_NAME", "lms_system")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Server configuration
	cfg.Server.Port = getEnv("PORT", "8080")

	// Logger configuration
	cfg.Logger.Level = getEnv("LOG_LEVEL", "info")
	cfg.Logger.Format = getEnv("LOG_FORMAT", "json")

	// Keycloak configuration
	cfg.Keycloak.Host = getEnv("KEYCLOAK_HOST", "http://localhost:8081")
	cfg.Keycloak.Realm = getEnv("KEYCLOAK_REALM", "lms")
	cfg.Keycloak.ClientID = getEnv("KEYCLOAK_CLIENT_ID", "lms-client")
	cfg.Keycloak.ClientSecret = getEnv("KEYCLOAK_CLIENT_SECRET", "")
	cfg.Keycloak.AdminUser = getEnv("KEYCLOAK_ADMIN_USER", "admin")
	cfg.Keycloak.AdminPass = getEnv("KEYCLOAK_ADMIN_PASS", "admin")

	// MinIO configuration
	cfg.MinIO.Endpoint = getEnv("MINIO_ENDPOINT", "localhost:9000")
	cfg.MinIO.AccessKeyID = getEnv("MINIO_ACCESS_KEY_ID", "minioadmin")
	cfg.MinIO.SecretAccessKey = getEnv("MINIO_SECRET_ACCESS_KEY", "minioadmin123")
	cfg.MinIO.UseSSL = getEnv("MINIO_USE_SSL", "false") == "true"
	cfg.MinIO.BucketName = getEnv("MINIO_BUCKET_NAME", "lms-files")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
