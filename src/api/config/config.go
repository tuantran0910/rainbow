package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
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

var (
	once sync.Once
	cfg  *Config
)

func LoadConfig() error {
	var err error
	once.Do(func() {
		os.Setenv("TZ", "UTC")

		// Load database connection configuration
		host := getEnv("DB_HOST", "localhost")
		user := getEnv("API_DB_USER", "rainbow")
		password := getEnv("API_DB_PASSWORD", "")
		if password == "" {
			err = fmt.Errorf("API_DB_PASSWORD is required but not set")
			return
		}

		dbName := getEnv("API_DB_NAME", "rainbow")
		dbPort := getEnv("DB_PORT", "5432")
		serverPort := getEnv("SERVER_PORT", "5000")

		maxIdleConns, convErr := strconv.Atoi(getEnv("MAX_IDLE_CONNS", "10"))
		if convErr != nil {
			err = fmt.Errorf("invalid MAX_IDLE_CONNS: %v", convErr)
			return
		}

		maxOpenConns, convErr := strconv.Atoi(getEnv("MAX_OPEN_CONNS", "20"))
		if err != nil {
			err = fmt.Errorf("invalid MAX_OPEN_CONNS: %v", convErr)
		}

		connMaxLifetime, convErr := time.ParseDuration(getEnv("CONN_MAX_LIFETIME", "1h"))
		if err != nil {
			err = fmt.Errorf("invalid CONN_MAX_LIFETIME: %v", convErr)
		}

		migrationsDir := getEnv("MIGRATIONS_DIR", "migrations")

		dbDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			host, user, password, dbName, dbPort)

		cfg = &Config{
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
		}
	})

	return err
}

func GetConfig() *Config {
	if cfg == nil {
		LoadConfig()
	}
	return cfg
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
