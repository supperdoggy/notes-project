package main

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Note struct {
	Id       int                    `form:"_id" json:"_id"`
	PublicId string                 `form:"publicId" json:"publicId"`
	Title    string                 `form:"title" json:"title"`
	Text     string                 `form:"text" json:"text"`
	Owner    string                 `form:"owner" json:"owner"`
	Created  time.Time              `form:"created" json:"created"`
	Shared   bool                   `form:"shared" json:"shared"`
	Users    []string `form:"users" json:"users"`
	gorm.Model
}

func (n *Note) addNewUser(userId string, p Permissions) error {
	if !n.Shared {
		return fmt.Errorf("note is not shared")
	}
	for _, v := range n.Users {
		if v == userId {
			return fmt.Errorf("user already in map")
		}
	}
	n.Users = append(n.Users, userId)
	return nil
}

func (n *Note) deleteUser(userId string) error {
	if !n.Shared {
		return fmt.Errorf("note is not shared")
	}
	for k, v := range n.Users{
		if v == userId {
			n.Users[k] = ""
		}
	}
	return nil
}
