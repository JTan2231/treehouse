package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func MakePost(c *gin.Context) {
	c.HTML(http.StatusOK, "blogPost.tmpl", gin.H{})
}