package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
)

type obj map[string]interface{}

func notePage(c *gin.Context) {
	id := c.Param("id")
	checkLogin(c)
	// getting token struct
	t, err := c.Cookie("t")
	if err != nil {
		c.Redirect(308, "/auth/login")
		return
	}
	token, err := findTokenStructInMap(t)
	if err != nil {
		c.Redirect(308, "/auth/login")
		return
	}
	result := make(map[string]interface{})
	resp, err := http.PostForm("http://localhost:2020/api/getNote", url.Values{"id": {id}, "t": {t}})
	if err != nil {
		log.Println("error:", err.Error())
		c.Redirect(308, "/")
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("error:", err.Error())
		c.Redirect(308, "/")
		return
	}
	if result["ok"].(bool) {
		note := result["answer"].(map[string]interface{})
		c.HTML(200, "comment.html", obj{"token": token, "note": note, "id": id})
		return
	} else {
		log.Println("not ok", result["error"])
		c.Redirect(308, "/")
	}
}
