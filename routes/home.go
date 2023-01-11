package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
)

func ServeHome(c *gin.Context) {
    session, _ := config.Store.Get(c.Request, "session")

	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"API_ROOT": config.API_ROOT,
        "username": session.Values["username"],
	})
}
