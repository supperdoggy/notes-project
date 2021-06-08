package main

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id       int       `form:"_id" json:"_id"`
	UniqueId int       `form:"uniqueId" json:"unique_id"`
	Name     string    `form:"name" json:"name"`
	Username string    `form:"username" json:"username"`
	Password string    `form:"password" json:"password"`
	Created  time.Time `form:"created" json:"created"`
	gorm.Model
}
