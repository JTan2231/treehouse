package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
)

func ServeHome(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")

	_, ok := session.Values["username"]
	if !ok {
		c.HTML(200, "404_redirect.tmpl", gin.H{
			"url": "/",
		})
	} else {
		c.HTML(http.StatusOK, "home.tmpl", gin.H{})
	}
}
