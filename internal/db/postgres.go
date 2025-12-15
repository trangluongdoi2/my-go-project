package db

import (
	"fmt"
	"go-backend-project/internal/booking"
	"go-backend-project/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseService struct {
	DB *gorm.DB
}

func NewPostgresConnection(cfg *config.Config) (*DatabaseService, error) {
	host := cfg.Database.Host
	port := cfg.Database.Port
	user := cfg.Database.Username
	password := cfg.Database.Password
	dbName := cfg.Database.DbName

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	pgConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
		WithoutReturning:     false,
	}

	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// This is connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	databaseService := DatabaseService{}
	databaseService.DB = db
	databaseService.InitTables()

	return &databaseService, nil
}

var dbModels = []interface{}{
	&booking.Booking{},
}

func (d *DatabaseService) InitTables() {
	err := d.DB.AutoMigrate(dbModels...)
	if err != nil {
		fmt.Println(err.Error(), "err when inittables...")
	}
}
