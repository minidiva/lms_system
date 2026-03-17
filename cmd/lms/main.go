package main

import (
	"fmt"
	"lms_system/internal/infrastructure/clients/http/keycloak"
	minioClient "lms_system/internal/infrastructure/clients/minio"
	"lms_system/internal/service/auth"
	"lms_system/internal/service/file"
	"lms_system/internal/service/lms"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lms_system/config"
	_ "lms_system/docs"
	httpDelivery "lms_system/internal/delivery/http"
	infraPostgres "lms_system/internal/infrastructure/repository/postgres"
)

// @title           LMS System API
// @version         1.0
// @description     API для системы управления обучением
// @host            localhost:8082
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load and parse config
	cfg := config.LoadConfig()

	// Initialize logger
	logger := initLogger(cfg)

	// Initialize database
	db, err := initDatabase(cfg.GetDatabaseDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	// Initialize repository
	mainRepo := infraPostgres.NewMainRepository(db, logger)

	// Initialize Keycloak client
	keycloakClient := keycloak.NewClient(
		cfg.Keycloak.Host,
		cfg.Keycloak.Realm,
		cfg.Keycloak.ClientID,
		cfg.Keycloak.ClientSecret,
		cfg.Keycloak.AdminUser,
		cfg.Keycloak.AdminPass,
	)

	// Initialize MinIO client
	minioClientInstance, err := minioClient.NewClient(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize MinIO client:", err)
	}

	// Initialize services
	fileSvc := file.NewService(minioClientInstance)
	svc := lms.NewService(mainRepo, logger, fileSvc)
	authSvc := auth.NewService(mainRepo, logger, keycloakClient)

	// Initialize HTTP server
	server := httpDelivery.NewServer(svc, authSvc, fileSvc, cfg.Server.Port)
	if err := server.Start(); err != nil {
		logger.Fatal("Server failed to start:", err)
	}
}

func initDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func initLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set log format
	if cfg.Logger.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return logger
}
