package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"

    config "treehouse/config"
)

func ServeLogin(c *gin.Context) {   
    c.HTML(http.StatusOK, "login.tmpl", gin.H{
        "API_ROOT": config.API_ROOT,
    })
}
