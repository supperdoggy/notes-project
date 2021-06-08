package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type cache struct {
	m map[string]enterToken
	sync.RWMutex
}

var tokenCache = cache{
	m:       make(map[string]enterToken),
	RWMutex: sync.RWMutex{},
}

var usersCollection, err = getDBSession()

func getBcrypt(text string) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(text), 4)
	if err != nil {
		panic(err.Error())
	}
	return string(hashedPass)
}

func main() {
	fmt.Println("Starting server...")
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/login", login)
		api.POST("/register", register)
		api.POST("/token", validateToken)
		api.POST("/newToken", newToken)
		api.POST("/getTokenStruct", getTokenStruct)
		api.POST("/deleteToken", deleteCookieFromMap)
		api.POST("/admin", userIsAdmin)
		api.POST("/deleteUser", deleteUser)
		api.POST("/getAllUsers", getAllUsers)
	}

	if err := r.Run(":2283"); err != nil {
		fmt.Println(err.Error())
		return
	}
}
