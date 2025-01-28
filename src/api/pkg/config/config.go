package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type DatabaseConfig struct {
	DatabaseDsn     string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type ServerConfig struct {
	MigrationsDir string
	ServerPort    string
}

type Config struct {
	DatabaseConfig *DatabaseConfig
	ServerConfig   *ServerConfig
}

func LoadConfig() (*Config, error) {
	os.Setenv("TZ", "UTC")

	// Load database connection configuration
	host := GetEnv("DB_HOST", "localhost")
	user := GetEnv("API_DB_USER", "rainbow")
	password := GetEnv("API_DB_PASSWORD", "R&inb0w2024!Data")
	dbName := GetEnv("API_DB_NAME", "rainbow")
	dbPort := GetEnv("DB_PORT", "5432")
	serverPort := GetEnv("SERVER_PORT", "5000")

	maxIdleConns, err := strconv.Atoi(GetEnv("MAX_IDLE_CONNS", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_IDLE_CONNS: %v", err)
	}

	maxOpenConns, err := strconv.Atoi(GetEnv("MAX_OPEN_CONNS", "20"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_OPEN_CONNS: %v", err)
	}

	connMaxLifetime, err := time.ParseDuration(GetEnv("CONN_MAX_LIFETIME", "1h"))
	if err != nil {
		return nil, fmt.Errorf("invalid CONN_MAX_LIFETIME: %v", err)
	}

	migrationsDir := GetEnv("MIGRATIONS_DIR", "migrations")

	dbDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbName, dbPort)

	return &Config{
		DatabaseConfig: &DatabaseConfig{
			DatabaseDsn:     dbDsn,
			MaxIdleConns:    maxIdleConns,
			MaxOpenConns:    maxOpenConns,
			ConnMaxLifetime: connMaxLifetime,
		},
		ServerConfig: &ServerConfig{
			MigrationsDir: migrationsDir,
			ServerPort:    serverPort,
		},
	}, nil
}

// Helper function to get environment variables with defaults
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
