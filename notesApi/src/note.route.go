package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func updateNote(c *gin.Context) {
	id := c.PostForm("id")
	Text := c.PostForm("Text")
	Title := c.PostForm("Title")

	if id == "" {
		c.JSON(200, map[string]interface{}{
			"ok":     false,
			"error":  "id is empty",
			"answer": false,
		})
		return
	}
	var n Note
	if err := notesSession.Where("public_id = ?", id).First(&n).Error; err != nil {
		c.JSON(400, map[string]interface{} {
			"ok": false,
			"error": "no note",
			"answer": false,
		})
		return
	}
	n.Text = Text
	n.Title = Title
	if err := notesSession.Save(&n).Error; err != nil {
		c.JSON(400, map[string]interface{} {
			"ok": false,
			"error": "error saving note",
			"answer": false,
		})
	}
	c.JSON(200, map[string]interface{}{
		"ok":     true,
		"error":  "",
		"answer": true,
	})
	return
}

func newNote(c *gin.Context) {
	Title := c.PostForm("Title")
	Text := c.PostForm("Text")
	Username := c.PostForm("Username")
	if Title == "" || Username == "" {
		c.JSON(200, map[string]interface{}{
			"ok":     false,
			"error":  "not all fields are filled",
			"answer": false,
		})
		return
	}
	note := Note{
		PublicId: strconv.FormatInt(time.Now().UnixNano(), 10), // just taking current nanosecs in unix
		Title:    Title,
		Text:     Text,
		Owner:    Username,
		Created:  time.Now(),
		Shared:   false,
		Users:    nil,
	}
	err = notesSession.Create(&note).Error
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"ok":     false,
			"error":  err.Error(),
			"answer": false,
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok":     true,
		"error":  "",
		"answer": true,
	})
	return
}

func getNote(c *gin.Context) {
	response := map[string]interface{}{
		"ok":     false,
		"error":  "",
		"answer": nil,
	}
	id := c.PostForm("id")
	t := c.PostForm("t")
	var username string

	var result Note
	err = notesSession.Model(&Note{}).Where("public_id = ?", id).First(&result).Error
	if err != nil {
		log.Println(err.Error())
		response["error"] = err.Error()
		c.JSON(400, response)
		return
	}

	// getting username
	data := make(map[string]interface{})
	resp, err := http.PostForm("http://localhost:2283/api/getTokenStruct", url.Values{"t": {t}})
	if err != nil {
		response["error"] = err.Error()
		c.JSON(200, response)
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		response["error"] = err.Error()
		c.JSON(400, response)
		return
	}
	// checking if user allowed to take get note
	if data["ok"].(bool) {
		username = data["answer"].(string)
	} else {
		response["error"] = "wrong request"
		c.JSON(400, response)
		return
	}
	if result.Owner == username {
		response["ok"] = true
		response["answer"] = result
		c.JSON(200, response)
		return
	} else if result.Shared {
		for _, v := range result.Users {
			if v == username {
				response["ok"] = true
				response["answer"] = result
				c.JSON(200, response)
				return			}
		}
	} else {
		response["error"] = "not allowed"
		c.JSON(400, response)
		return
	}
}

func shareNote(c *gin.Context) {
	owner := c.PostForm("Owner")
	username := c.PostForm("Username") // username of user we want to share note with
	id := c.PostForm("Id")             // public id
	var note Note

	err = notesSession.Where("public_id = ?", id).First(&note).Error
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"ok":     false,
			"error":  err.Error(),
			"answer": false,
		})
		return
	}

	note.shareNote()
	err = note.addNewUser(username, Permissions{
		CanRedact:      true,
		CanAddNewUsers: true,
	})
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"ok":     false,
			"error":  err.Error(),
			"answer": false,
		})
		return
	}

	var n Note
	if err := notesSession.Where("public_id = ? AND Owner = ?", id, owner).First(&n).Error; err != nil {
		c.JSON(400, map[string]interface{} {
			"ok": false,
			"error": "no note",
			"answer": false,
		})
		return
	}
	n.Shared = true
	n.Users = note.Users
	err = notesSession.Model(&Note{}).Exec("UPDATE \"notes\" SET \"users\"=?::text[],\"shared\"=true,\"updated_at\"=? WHERE id = ?", n.Users, n.UpdatedAt, n.Id).Error
	//if err := notesSession.Model(&Note{}).Where("id = ?", n.ID).Updates(map[string]interface{} {"shared":true, "Users":gorm.Expr("?::text[]", note.Users)}).Error; err != nil {
	if err != nil {
		c.JSON(400, map[string]interface{} {
			"ok": false,
			"error": "error saving note",
			"answer": false,
		})
	}

	c.JSON(200, map[string]interface{}{
		"ok":     true,
		"error":  "",
		"answer": true,
	})
	return
}

func sendNotes(c *gin.Context) {
	username := c.PostForm("username")
	var ownedNotes []Note
	err = notesSession.Where("Owner = ?", username).Find(&ownedNotes).Error
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"ok":     false,
			"error":  err.Error(),
			"answer": "",
		})
		return
	}
	var sharedNotes []Note
	var allNotes []Note
	err = notesSession.Find(&allNotes).Error
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"ok":     false,
			"error":  err.Error(),
			"answer": false,
		})
		return
	}
	for _, v := range allNotes {
		if v.Shared{
			for _, j := range v.Users {
				if j == username {
					sharedNotes = append(sharedNotes, v)
				}
			}
		}
	}
	c.JSON(200, map[string]interface{}{
		"ownedNotes":  ownedNotes,
		"sharedNotes": sharedNotes,
	})
	return
}

func deleteNote(c *gin.Context) {
	response := map[string]interface{}{
		"ok":    false,
		"error": "",
	}
	id := c.PostForm("id")

	err = notesSession.Model(&Note{}).Find(id).Delete(&Note{}).Error
	if err != nil {
		response["error"] = err.Error()
		c.JSON(400, response)
		return
	}
	response["ok"] = true
	c.JSON(200, response)
	return
}
