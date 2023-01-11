package main

import (
	"github.com/gin-gonic/gin"
	"treehouse/config"
	"treehouse/db"
	"treehouse/routes"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

func main() {
	db.InitDB()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./styles")

	router.GET("/", routes.ServeLogin) 
	router.GET("/newuser", routes.ServeNewUser)

    router.GET("/:username", routes.ServeProfile) 
    router.GET("/:username/:slug", routes.GetArticle)

	router.POST("/login", routes.AuthenticateLogin)
	router.POST("/newuser", routes.CreateNewUser)

	authRouter := router.Group("/create", routes.AuthRequired)
	authRouter.GET("/create-article", routes.GetCreateArticle)
	authRouter.POST("/articles", routes.CreateArticle)

	router.Run(config.DOMAIN)
}
