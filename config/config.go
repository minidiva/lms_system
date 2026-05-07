package config

import "fmt"

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Server struct {
		Port string
	}
	Logger struct {
		Level  string
		Format string
	}
	Keycloak struct {
		Host         string
		Realm        string
		ClientID     string
		ClientSecret string
		AdminUser    string
		AdminPass    string
	}
	MinIO struct {
		Endpoint        string
		PublicEndpoint  string // ← внешний адрес для presigned URL
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
		BucketName      string
	}
}

func (cfg *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
}
