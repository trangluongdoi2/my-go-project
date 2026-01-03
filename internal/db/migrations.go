package db

import (
	"fmt"

	"gorm.io/gorm"
)

type Migration struct {
	ID      string
	Execute func(*gorm.DB) error
}

func (d *DatabaseService) runMigrations() error {
	migrations := []Migration{}

	for _, m := range migrations {
		if err := m.Execute(d.DB); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.ID, err)
		}
		fmt.Printf("Migration %s completed\n", m.ID)
	}

	return nil
}
