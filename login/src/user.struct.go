package main

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id       int `form:"_id" json:"_id"`
	UniqueId string        `form:"uniqueId" json:"unique_id"`
	Username string        `form:"username" json:"username"`
	Password string        `form:"password" json:"password"`
	Created  time.Time     `form:"created" json:"created"`
	IsAdmin  bool          `form:"isAdmin" json:"isAdmin"`
	gorm.Model
}
