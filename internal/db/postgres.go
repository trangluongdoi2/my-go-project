package db

import (
	"fmt"
	"go-backend-project/internal/config"
	"log"
	"strings"
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

	// This is CONNECTION POOL
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	databaseService := DatabaseService{}
	databaseService.DB = db
	databaseService.InitTables()

	return &databaseService, nil
}

var dbModels = []interface{}{
	// &appointment.Appointment{},
	// &staff.Staff{},
	// &serviceoffering.ServiceOffering{},
	// &serviceoffering.StaffService{},
	// &user.User{},
}

func (d *DatabaseService) SafeAutoMigrate(models ...interface{}) error {
	for _, model := range models {
		db := d.DB
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return fmt.Errorf("failed to parse model %T: %w", model, err)
		}

		tableName := stmt.Schema.Table
		log.Printf("[SafeAutoMigrate] Processing table: %s", tableName)

		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("auto migrate failed for %T: %w", model, err)
		}
		log.Printf("[SafeAutoMigrate] AutoMigrate completed for %s", tableName)

		expectedIndexes := make(map[string]bool)
		for _, idx := range stmt.Schema.ParseIndexes() {
			if idx == nil || len(idx.Fields) == 0 {
				continue
			}
			expectedIndexes[idx.Name] = (idx.Class == "UNIQUE")
		}

		var existingCols []string
		colQuery := `
			SELECT column_name
			FROM information_schema.columns
			WHERE table_name = $1
		`
		if err := db.Raw(colQuery, tableName).Scan(&existingCols).Error; err != nil {
			return fmt.Errorf("failed fetching columns for %s: %w", tableName, err)
		}

		expectedCols := make(map[string]bool)
		for _, field := range stmt.Schema.Fields {
			expectedCols[field.DBName] = true
		}

		existingColsMap := make(map[string]bool)
		for _, col := range existingCols {
			existingColsMap[col] = true
		}
		for _, col := range existingCols {
			if _, ok := expectedCols[col]; !ok {
				log.Printf("[SafeAutoMigrate] Dropping unused column %s on table %s", col, tableName)
				if err := db.Migrator().DropColumn(model, col); err != nil {
					return fmt.Errorf("failed dropping column %s on %s: %w", col, tableName, err)
				}
			}
		}

		var existing []struct {
			Name   string
			Unique bool
		}
		query := `
			SELECT ic.relname AS name, i.indisunique
			FROM pg_index i
			JOIN pg_class t ON t.oid = i.indrelid
			JOIN pg_class ic ON ic.oid = i.indexrelid
			WHERE t.relname = ?
		`
		if err := db.Raw(query, tableName).Scan(&existing).Error; err != nil {
			return fmt.Errorf("failed fetching indexes for %s: %w", tableName, err)
		}

		for _, ex := range existing {
			if strings.HasSuffix(ex.Name, "_pkey") || strings.HasPrefix(ex.Name, "pg_") {
				continue
			}

			expectedUnique, ok := expectedIndexes[ex.Name]

			if !ok {
				log.Printf("[SafeAutoMigrate] Dropping unused index %s on table %s", ex.Name, tableName)
				if err := db.Migrator().DropIndex(model, ex.Name); err != nil {
					return fmt.Errorf("failed dropping unused index %s on %s: %w", ex.Name, tableName, err)
				}
				continue
			}

			if ex.Unique != expectedUnique {
				log.Printf("[SafeAutoMigrate] Dropping conflicting index %s on table %s (expected unique=%v, got unique=%v)",
					ex.Name, tableName, expectedUnique, ex.Unique)

				if err := db.Migrator().DropIndex(model, ex.Name); err != nil {
					return fmt.Errorf("failed dropping index %s on %s: %w", ex.Name, tableName, err)
				}
			}
		}

		for idxName := range expectedIndexes {
			if !db.Migrator().HasIndex(model, idxName) {
				log.Printf("[SafeAutoMigrate] Creating index %s on table %s", idxName, tableName)
				if err := db.Migrator().CreateIndex(model, idxName); err != nil {
					return fmt.Errorf("failed creating index %s on %s: %w", idxName, tableName, err)
				}
			}
		}
	}
	log.Println("[SafeAutoMigrate] Migration completed successfully")
	return nil
}

func (d *DatabaseService) InitTables() {
	fmt.Println("AUTO-MIGRATION TABLES...")
	if err := d.SafeAutoMigrate(dbModels...); err != nil {
		fmt.Printf("AUTO-MIGRATION [ERROR]: %v\n", err)
		return
	}

	fmt.Println("SCHEMA SYNC COMPLETED")
}
