package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func adminPage(c *gin.Context) {
	checkLogin(c)
	checkIfAdmin(c)

	t, err := c.Cookie("t")
	if err != nil {
		fmt.Println("cookie error")
		c.Redirect(308, "/auth/login")
		return
	}
	resp, err := http.PostForm("http://localhost:2283/api/getAllUsers", url.Values{"t": {t}})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	answer := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&answer)
	if err != nil {
		fmt.Println("Decode error")
		c.Redirect(308, "/")
		return
	}
	u := make([]interface{}, 100)
	if answer["ok"].(bool) {
		u = answer["answer"].([]interface{})
	} else {
		fmt.Println("returning")
		return
	}
	users := processAnswer(u)
	fmt.Println(users[1])
	c.HTML(200, "adminPage.html", obj{"users": users[1:]})
	return
}
