package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/farmako/coupon-system/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConn represents a database connection
type DBConn struct {
	DB *gorm.DB
}

// Config represents database configuration
type Config struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection creates a new database connection
func NewConnection(config Config) (*DBConn, error) {
	var db *gorm.DB
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	switch config.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open("coupon_system.db"), &gorm.Config{
			Logger: newLogger,
		})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	if err != nil {
		return nil, err
	}

	return &DBConn{DB: db}, nil
}

// Migrate performs database migrations
func (d *DBConn) Migrate() error {
	return d.DB.AutoMigrate(&models.Coupon{}, &models.CouponUsage{})
}

// Close closes the database connection
func (d *DBConn) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
