package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
)

func ServeLanding(c *gin.Context) {
    //session, _ := config.Store.Get(c.Request, "session")

    //_, ok := session.Values["username"]
    if false {
        c.HTML(200, "404_redirect.tmpl", gin.H{
            "url": "/home",
        })
    } else {
        c.HTML(http.StatusOK, "landing.tmpl", gin.H{
            "API_ROOT": config.API_ROOT,
        })
    }
}
