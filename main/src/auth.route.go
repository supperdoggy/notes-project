package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func loginPage(c *gin.Context) {
	cookie, err := c.Cookie("error")
	if err != nil {
		cookie = ""
	}
	data := gin.H{}
	if cookie != "" {
		data["error"] = cookie
	}

	c.HTML(200, "login.html", data)
}

func registerPage(c *gin.Context) {
	cookie, err := c.Cookie("error")
	if err != nil {
		cookie = ""
	}
	data := gin.H{}
	if cookie != "" {
		data["error"] = cookie
	}

	c.HTML(200, "register.html", data)
}

func login(c *gin.Context) {
	// deleting token
	deleteCookieFromMap(c)
	c.SetCookie("t", "", -1, "/auth/login", "localhost", false, true)

	username := c.PostForm("login")
	password := c.PostForm("pass")

	if username == "" || password == "" {
		c.SetCookie("error", "Fill all the inputs!", 20, "/auth/login", "localhost", false, true)
		c.Redirect(http.StatusMovedPermanently, "/auth/login")
		return
	}

	if validateUser(username, password) == true {
		createNewTokenCookie(c, username)
		c.SetCookie("error", "", -1, "/auth/login", "localhost", false, true)
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	c.SetCookie("error", "Wrong username/password", 20, "/auth/login", "localhost", false, true)
	c.Redirect(http.StatusMovedPermanently, "/auth/login")
}

func validateUser(username, password string) interface{} {
	resp, err := http.PostForm("http://localhost:2283/api/login", url.Values{"login": {username}, "pass": {password}})
	if err != nil {
		return false
	}

	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false
	}
	if data["ok"] == true {
		return data["answer"]
	} else {
		fmt.Println(data["error"])
	}
	return false
}

func register(c *gin.Context) {
	// deleting token
	deleteCookieFromMap(c)
	c.SetCookie("t", "", -1, "/auth/login", "localhost", false, true)

	username := c.PostForm("login")
	password := c.PostForm("pass")
	if username == "" || password == "" {
		c.SetCookie("error", "Fill all the inputs", 20, "/auth/register", "localhost", false, true)
		c.Redirect(http.StatusMovedPermanently, "/auth/register")
		return
	}

	resp, err := http.PostForm("http://localhost:2283/api/register", url.Values{"login": {username}, "pass": {password}})
	if err != nil {
		panic(err.Error())
	}

	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err.Error())
	}
	if data["ok"] == true {
		c.Redirect(308, "/auth/login")
		return
	} else {
		fmt.Println()
		c.SetCookie("error", fmt.Sprintf("%v", data["error"]), 20, "/auth/register", "localhost", false, true)
		c.Redirect(308, "/auth/register")
		return
	}
}

func checkLogin(c *gin.Context) {
	t, err := c.Cookie("t")
	if err != nil || !validateEntryToken(&t) {
		c.Redirect(http.StatusPermanentRedirect, "/auth/login")
		return
	}
}

func checkIfAdmin(c *gin.Context) {
	t, err := c.Cookie("t")
	if err != nil || !userIsAdmin(&t) {
		c.Redirect(308, "/auth/login")
		return
	}
}

func userIsAdmin(s *string) bool {
	resp, err := http.PostForm("http://localhost:2283/api/admin", url.Values{"t": {*s}})
	if err != nil {
		return false
	}
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println(data)
	if !data["ok"].(bool) {
		return false
	} else {
		return data["answer"].(bool)
	}
}
