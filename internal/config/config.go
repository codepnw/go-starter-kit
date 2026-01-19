package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/codepnw/go-starter-kit/pkg/utils/validate"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	APP AppConfig `envPrefix:"APP_"`
	DB  DBConfig  `envPrefix:"DB_"`
}

type AppConfig struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"8080"`
}

type DBConfig struct {
	User     string `env:"USER" validate:"required"`
	Password string `env:"PASSWORD" validate:"required"`
	Name     string `env:"NAME" validate:"required"`
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"5432"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

func LoadConfig(path string) (*EnvConfig, error) {
	// Load .env file
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("load env failed: %w", err)
	}

	cfg := new(EnvConfig)
	// Parse env tags
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env failed: %w", err)
	}
	
	// Validate config
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate env failed: %w", err)
	}
	return cfg, nil
}

// GetDatabaseDSN : Default PostgreSQL
func (cfg *EnvConfig) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)
}

func (cfg *EnvConfig) GetAppAddress() string {
	return fmt.Sprintf(":%d", cfg.APP.Port)
}