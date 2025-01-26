package postgres

import (
	"fmt"
	"sync"

	"github.com/tuantran0910/rainbow/internal/models/product"
	"github.com/tuantran0910/rainbow/pkg/config"
	"github.com/tuantran0910/rainbow/pkg/utils/logger"
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

func (pc *PostgresConnector) migrateDatabaseTables(models ...interface{}) error {
	if err := pc.ins.AutoMigrate(models...); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (pc *PostgresConnector) GetInstance() (*gorm.DB, error) {
	var err error
	pc.once.Do(func() {
		if pc.ins == nil {
			// Validate configs
			if err := pc.validateConfig(); err != nil {
				log.Error(err.Error())
				return
			}

			// Make a connection to a database
			var db *gorm.DB
			db, err = gorm.Open(postgres.Open(pc.cfg.DatabaseConfig.DatabaseDsn), &gorm.Config{})
			if err != nil {
				log.Error(err.Error())
				return
			}

			// Set the instance
			pc.ins = db

			// Set pool configuration
			if err := pc.configureConnectionPool(); err != nil {
				pc.ins = nil
				return
			}

			// Auto table migration
			if err := pc.migrateDatabaseTables(&product.Product{}); err != nil {
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
			log.Error(err.Error())
			return err
		}

		// Close the database connection
		if err := db.Close(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}
