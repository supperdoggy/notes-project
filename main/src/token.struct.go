package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

type enterToken struct {
	Token     string `bson:"token" json:"token" form:"token"`
	Limited   bool   `bson:"limited" json:"limited" form:"limited"`
	SavedTime int64  `bson:"savedTime" json:"saved_time" form:"savedTime"`
	Username  string `bson:"username" json:"username"`
}

func (t *enterToken) expired(minutes int64) (result bool) {
	if t.Limited == true {
		result = !(((time.Now().Unix() - t.SavedTime) / 60) > minutes)
	}
	return false
}

func createNewToken(username string) (string, error) {
	request := url.Values{"username": {username}}
	resp, err := http.PostForm("http://localhost:2283/api/newToken", request)
	if err != nil {
		return "", fmt.Errorf("error requesting new token")
	}
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if !data["ok"].(bool) {
		return "", data["error"].(error)
	}
	return data["answer"].(string), nil
}

func createNewTokenCookie(c *gin.Context, username string) {
	t, err := createNewToken(username) // cookie is limited
	if err != nil {
		c.Redirect(308, "/auth/login")
		return
	}
	c.SetCookie("t", t, 999, "/", "localhost", false, true)
}

func validateEntryToken(s *string) bool {
	resp, err := http.PostForm("http://localhost:2283/api/token", url.Values{"t": {*s}})
	if err != nil {
		fmt.Println("http://localhost:2283/api/newToken isn't responding")
		return false
	}
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("error decoding json data from http://localhost:2283/api/newToken")
		return false
	}
	if !data["ok"].(bool) {
		return false
	} else {
		return data["answer"].(bool)
	}
}

func findTokenStructInMap(t string) (string, error) {
	resp, err := http.PostForm("http://localhost:2283/api/getTokenStruct", url.Values{"t": {t}})
	if err != nil {
		return "", err
	}
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	if !data["ok"].(bool) {
		return "", fmt.Errorf(data["error"].(string))
	} else {
		return data["answer"].(string), nil
	}
}

func deleteCookieFromMap(c *gin.Context) {
	t, err := c.Cookie("t")
	if err != nil {
		return
	}
	_, _ = http.PostForm("http://localhost:2283/api/deleteToken", url.Values{"t": {t}})
	return
}
