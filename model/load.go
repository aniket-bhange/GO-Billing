package model

import "github.com/jinzhu/gorm"

func LoadModels(db *gorm.DB) *gorm.DB {
	db.Debug().AutoMigrate(&Users{}, &Client{}, &Product{})
	return db
}
