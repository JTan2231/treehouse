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
        c.HTML(http.StatusOK, "landing.tmpl", gin.H{
            "API_ROOT": config.API_ROOT,
        })
    } else {
        c.HTML(http.StatusOK, "home.tmpl", gin.H{})
    }
}
