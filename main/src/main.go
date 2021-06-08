package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server...")
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	auth := r.Group("/auth")
	{
		// login
		auth.GET("/login", loginPage)
		auth.POST("/login", login)
		// register
		auth.GET("/register", registerPage)
		auth.POST("/register", register)
	}

	m := r.Group("/")
	{
		m.POST("/", mainPage)
		m.GET("/", mainPage)
		m.Any("/note/:id", notePage)
		m.GET("/admin", adminPage)
	}

	if err := r.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
