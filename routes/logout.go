package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
)

func HandleLogout(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)

	c.HTML(http.StatusOK, "404_redirect.tmpl", gin.H{
		"url": "/",
	})
}
