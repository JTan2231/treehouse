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
	router.Static("/assets", "./assets")

	router.GET("/", routes.ServeLanding)
	router.GET("/home", routes.ServeHome)

	router.GET("/login", routes.ServeLogin)
	router.GET("/newuser", routes.ServeNewUser)

	router.GET("/:username", routes.ServeProfile)
	router.GET("/:username/:slug", routes.GetArticle)

	router.POST("/login", routes.AuthenticateLogin)
	router.GET("/logout", routes.HandleLogout)
	router.POST("/newuser", routes.CreateNewUser)

	router.POST("/comments", routes.CreateComment)
	router.GET("/comments", routes.GetComments)

	createGroup := router.Group("/create", routes.AuthRequired)
	createGroup.GET("/create-article", routes.GetCreateArticle)
	createGroup.POST("/create-article", routes.CreateArticle)

	authRouter := router.Group("/", routes.AuthRequired)
	authRouter.POST("/subscribe", routes.SubscribeToUser)
	authRouter.POST("/favorite", routes.FavoriteArticle)
	authRouter.POST("/articles", routes.CreateArticle)

	router.Run(config.DOMAIN)
}
