package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func getMongoSession() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=mmarchy password=abc123 dbname=notesnotes port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil || db == nil {
		return nil, err
	}
	if err = db.AutoMigrate(&Note{}); err != nil {
		log.Println("AutoMigrate() -", err.Error())
	}
	if err = db.AutoMigrate(&Permissions{}); err != nil {
		log.Println("AutoMigrate() -", err.Error())
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		log.Println("AutoMigrate() -", err.Error())
	}
	db.Model(&Note{})
	return db, nil
}