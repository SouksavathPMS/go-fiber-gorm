package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `gorm:"primaryKey;autoIncrement:true"`
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Publisher *string `json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
