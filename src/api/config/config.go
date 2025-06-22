package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	once sync.Once
	cfg  *Config
)

type DatabaseConfig struct {
	DatabaseDsn     string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type ServerConfig struct {
	MigrationsDir  string
	ServerPort     string
	ServiceName    string
	ServiceVersion string
	JWTSecret      string
}

type LoggerConfig struct {
	LogLevel      string
	EnableConsole bool
}

type Config struct {
	Environment    string
	DatabaseConfig *DatabaseConfig
	ServerConfig   *ServerConfig
	LoggerConfig   *LoggerConfig
}

func loadDatabaseConfig(env string) (*DatabaseConfig, error) {
	// Load database connection configuration
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "rainbow")
	password := getEnv("DB_PASSWORD", "")
	if password == "" {
		return nil, fmt.Errorf(
			"DB_PASSWORD is required but not set",
		)
	}

	dbName := getEnv("DB_NAME", "rainbow")
	dbPort := getEnv("DB_PORT", "5432")

	// Determine SSL mode based on environment and USE_CLOUD_SQL flag
	useCloudSQL := getEnv("USE_CLOUD_SQL", "false")
	var sslMode string
	if useCloudSQL == "true" || env == "production" {
		sslMode = "require"
	} else {
		sslMode = "disable"
	}

	maxIdleConns, err := strconv.Atoi(getEnv("MAX_IDLE_CONNS", "10"))
	if err != nil {
		return nil, err
	}

	maxOpenConns, err := strconv.Atoi(getEnv("MAX_OPEN_CONNS", "20"))
	if err != nil {
		return nil, err
	}

	connMaxLifetime, err := time.ParseDuration(getEnv("CONN_MAX_LIFETIME", "1h"))
	if err != nil {
		return nil, err
	}

	// Construct the database DSN
	dbDsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host,
		user,
		password,
		dbName,
		dbPort,
		sslMode,
	)

	return &DatabaseConfig{
		DatabaseDsn:     dbDsn,
		MaxIdleConns:    maxIdleConns,
		MaxOpenConns:    maxOpenConns,
		ConnMaxLifetime: connMaxLifetime,
	}, nil
}

func loadServerConfig() (*ServerConfig, error) {
	serverPort := getEnv("SERVER_PORT", "5000")
	migrationsDir := getEnv("MIGRATIONS_DIR", "migrations")
	serverName := getEnv("SERVICE_NAME", "rainbow-api")
	serverVersion := getEnv("SERVICE_VERSION", "1.0.0")
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters long for security")
	}

	return &ServerConfig{
		MigrationsDir:  migrationsDir,
		ServerPort:     serverPort,
		ServiceName:    serverName,
		ServiceVersion: serverVersion,
		JWTSecret:      jwtSecret,
	}, nil
}

func loadLoggerConfig() (*LoggerConfig, error) {
	logLevel := getEnv("LOG_LEVEL", "debug")
	enableConsole, err := strconv.ParseBool(getEnv("ENABLE_CONSOLE", "true"))
	if err != nil {
		return nil, err
	}

	return &LoggerConfig{
		LogLevel:      logLevel,
		EnableConsole: enableConsole,
	}, nil
}

func LoadConfig() error {
	var err error
	once.Do(func() {
		os.Setenv("TZ", "UTC")

		// Get the environment
		env := getEnv("ENV", "development")
		if env != "production" && env != "development" {
			err = fmt.Errorf("invalid environment: %s", env)
			return
		}

		// Load database connection configuration
		databaseConfig, dbErr := loadDatabaseConfig(env)
		if dbErr != nil {
			err = fmt.Errorf("failed to load database configuration: %w", dbErr)
			return
		}

		// Load server configuration
		serverConfig, svErr := loadServerConfig()
		if svErr != nil {
			err = fmt.Errorf("failed to load server configuration: %v", svErr)
			return
		}

		// Get the logger configuration
		loggerConfig, lgErr := loadLoggerConfig()
		if lgErr != nil {
			err = fmt.Errorf("failed to load logger configuration: %v", lgErr)
			return
		}

		cfg = &Config{
			Environment:    env,
			DatabaseConfig: databaseConfig,
			ServerConfig:   serverConfig,
			LoggerConfig:   loggerConfig,
		}
	})

	return err
}

func GetConfig() (*Config, error) {
	if cfg == nil {
		if err := LoadConfig(); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
