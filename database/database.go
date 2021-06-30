package database

import (
	"e_commerce_furniture_with_fiber/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() *gorm.DB{
	dsn := "host=localhost user=postgres password=suat dbname=furniture port=5432 sslmode=disable TimeZone=Europe/Istanbul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&entity.Product{},&entity.Category{},&entity.User{},&entity.Faq{})
	if err != nil {
		return nil
	}
	return db
}