package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func getDBSession() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=mmarchy password=abc123 dbname=usersnotes port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil || db == nil {
		panic(err.Error())
		return nil, err
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		log.Println("AutoMigrate() -", err.Error())
	}
	return db, nil
}