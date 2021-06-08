package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func login(c *gin.Context) {
	response := map[string]interface{}{
		"ok":     false,
		"error":  "",
		"answer": false,
	}

	username := c.PostForm("login")
	password := c.PostForm("pass")
	if username == "" || password == "" {
		response["error"] = "not all inputs are filled"
		c.JSON(200, response)
		return
	}
	if validateUser(username, password, usersCollection) {
		response["ok"] = true
		response["answer"] = true
		c.JSON(200, response)
		return
	}
	response["ok"] = true
	c.JSON(200, response)
	return
}

func validateUser(username, password string, users *gorm.DB) bool {
	foundUsers := User{}
	err := users.Where("Username = ?", username).First(&foundUsers).Error
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundUsers.Password), []byte(password)); err == nil {
		return true
	}
	return false
}

func register(c *gin.Context) {
	response := map[string]interface{}{
		"ok":     false,
		"error":  "",
		"answer": false,
	}

	username := c.PostForm("login")
	password := c.PostForm("pass")
	fmt.Println(username, password)
	if username == "" || password == "" {
		c.JSON(200, map[string]interface{}{})
		return
	}
	u := User{
		UniqueId: strconv.FormatInt(time.Now().UnixNano(), 10),
		Username: username,
		Password: getBcrypt(password),
		Created:  time.Now(),
	}

	taken, err := usernameIsTaken(usersCollection, username)
	if taken {
		response["error"] = "Username is already taken!"
		c.JSON(200, response)
		return
	}
	err = usersCollection.Create(&u).Error
	if err != nil {
		response["error"] = "error inserting new user"
		c.JSON(200, response)
		return
	}
	response["ok"] = true
	response["answer"] = true
	c.JSON(200, response)
	return
}

func usernameIsTaken(users *gorm.DB, username string) (result bool, err error) {
	foundUsers := []User{}
	err = users.Where("Username = ?", username).First(&foundUsers).Error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// if user already in base
	if len(foundUsers) > 0 {
		result = true
		return
	}
	return
}

func validateToken(c *gin.Context) {
	t := c.PostForm("t")
	response := map[string]interface{}{
		"ok":     true,
		"error":  "",
		"answer": false,
	}
	// if token is not valid then delete it from cache
	if validateEntryToken(&t) {
		response["answer"] = true
	} else {
		response["answer"] = false
		tokenCache.Lock()
		defer tokenCache.Unlock()
		if _, ok := tokenCache.m[t]; ok {
			delete(tokenCache.m, t)
		}
	}
	c.JSON(200, response)
	return
}

func newToken(c *gin.Context) {
	response := map[string]interface{}{
		"ok":     false,
		"error":  "",
		"answer": false,
	}
	username := c.PostForm("username")
	if username == "" {
		response["error"] = "not provided username"
		c.JSON(200, response)
		return
	}
	t := createNewToken(true, username)
	tokenCache.Lock()
	defer tokenCache.Unlock()
	tokenCache.m[t.Token] = t
	response["ok"] = true
	response["answer"] = t.Token
	c.JSON(200, response)
	return
}

// takes token string and returns token struct
func getTokenStruct(c *gin.Context) {
	response := map[string]interface{}{
		"ok":     false,
		"error":  "",
		"answer": false,
	}
	t := c.PostForm("t")

	token, err := findTokenStructInMap(t)
	if err != nil {
		response["answer"] = err.Error()
		c.JSON(200, response)
		return
	} else {
		response["ok"] = true
		response["answer"] = token.Username
		c.JSON(200, response)
		return
	}
}

func userIsAdmin(c *gin.Context) {
	t := c.PostForm("t")
	response := map[string]interface{}{
		"ok":     false,
		"error":  "",
		"answer": false,
	}

	if err != nil {
		response["error"] = err.Error()
		c.JSON(200, response)
		return
	}
	tokenCache.Lock()
	defer tokenCache.Unlock()
	if token, ok := tokenCache.m[t]; ok {
		if !token.expired(30) {
			var u User
			err = usersCollection.Where("Username = ?", token.Username).First(&u).Error
			if err != nil {
				response["error"] = err.Error()
			} else {
				response["ok"] = true
				response["answer"] = u.IsAdmin
			}
		} else {
			delete(tokenCache.m, t)
		}
	}
	c.JSON(200, response)
	return
}

func deleteUser(c *gin.Context) {
	response := map[string]interface{}{
		"ok":     false,
		"error":  "",
		"answer": false,
	}
	id := c.PostForm("id")
	fmt.Println(id)
	err = usersCollection.Find(id).Delete(&User{}).Error
	if err != nil {
		response["error"] = err.Error()
		c.JSON(200, response)
		return
	}
	response["ok"] = true
	c.JSON(200, response)
	return
}

func getAllUsers(c *gin.Context) {
	response := map[string]interface{}{
		"ok":     false,
		"error":  "",
		"answer": false,
	}
	t := c.PostForm("t")
	if !validateEntryToken(&t) {
		response["error"] = "token error"
		c.JSON(200, response)
		return
	}
	users := make([]User, 100)
	if err != nil {
		response["error"] = err.Error()
		c.JSON(200, response)
		return
	}
	err = usersCollection.Find(&users).Error
	if err != nil {
		response["error"] = err.Error()
		c.JSON(200, response)
		return
	}
	response["answer"] = users
	response["ok"] = true
	c.JSON(200, response)
	return
}
