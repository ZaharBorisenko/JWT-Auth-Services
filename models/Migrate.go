package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}
