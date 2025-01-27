package postgres

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pressly/goose"
	"github.com/tuantran0910/rainbow/pkg/config"
	"github.com/tuantran0910/rainbow/pkg/utils/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = logger.GetLogger()

type IPostgresConnector interface {
	GetInstance() (*gorm.DB, error)
	Close() error
}

type PostgresConnector struct {
	ins  *gorm.DB
	cfg  *config.Config
	once sync.Once
}

func NewPostgresConnector(cfg *config.Config) IPostgresConnector {
	return &PostgresConnector{
		cfg: cfg,
	}
}

func (pc *PostgresConnector) validateConfig() error {
	if pc.cfg.DatabaseConfig == nil {
		return fmt.Errorf("invalid database configuration")
	}

	return nil
}

func (pc *PostgresConnector) configureConnectionPool() error {
	// Configure the database connection pool
	db, err := pc.ins.DB()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	db.SetMaxIdleConns(pc.cfg.DatabaseConfig.MaxIdleConns)
	db.SetMaxOpenConns(pc.cfg.DatabaseConfig.MaxOpenConns)
	db.SetConnMaxLifetime(pc.cfg.DatabaseConfig.ConnMaxLifetime)

	return nil
}

func (pc *PostgresConnector) GetInstance() (*gorm.DB, error) {
	var err error
	pc.once.Do(func() {
		if pc.ins == nil {
			// Validate configs
			if err := pc.validateConfig(); err != nil {
				log.Error("Failed to validate the database configuration %v", zap.Error(err))
				return
			}

			// Make a connection to a database
			var db *gorm.DB
			db, err = gorm.Open(postgres.Open(pc.cfg.DatabaseConfig.DatabaseDsn), &gorm.Config{})
			if err != nil {
				log.Error("Failed to connect to the database %v", zap.Error(err))
				return
			}

			// Convert GORM DB to Goose DB
			sqlDB, err := db.DB()
			if err != nil {
				log.Error("Failed to convert GORM DB to Goose DB %v", zap.Error(err))
				return
			}

			// Get the directory of the current file
			_, filePath, _, ok := runtime.Caller(0)
			if !ok {
				log.Error("Failed to get the current file path")
			}

			// Apply migrations
			migrationsDir := filepath.Join(filepath.Dir(filePath), "../migrations")
			if err := goose.Up(sqlDB, migrationsDir); err != nil {
				log.Error("Failed to apply migrations %v", zap.Error(err))
			}

			// Set the instance
			pc.ins = db

			// Set pool configuration
			if err := pc.configureConnectionPool(); err != nil {
				pc.ins = nil
				return
			}
		}
	})

	return pc.ins, nil
}

func (pc *PostgresConnector) Close() error {
	if pc.ins != nil {
		db, err := pc.ins.DB()
		if err != nil {
			log.Error("Failed to get the database connection %v", zap.Error(err))
			return err
		}

		// Close the database connection
		if err := db.Close(); err != nil {
			log.Error("Failed to close the database connection %v", zap.Error(err))
			return err
		}
	}

	return nil
}
