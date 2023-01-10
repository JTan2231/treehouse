package main

import (
    "github.com/gin-gonic/gin"
    "treehouse/db"
    "treehouse/routes"
    "treehouse/config"
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

    router.GET("/", routes.ServeLogin) // TODO: get an actual homepage
    router.GET("/newuser" ,routes.ServeNewUser)
    router.POST("/login", routes.AuthenticateLogin)


    authRouter := router.Group("/users/:username", routes.AuthRequired)
    authRouter.GET("/create-article",routes.GetCreateArticle)
    //authRouter.POST("/articles", routes.CreateArticle)
    authRouter.GET("/title", routes.GetArticle)
    //authRouter.GET("/users/:username/profile", routes.GetProfile)

    router.Run(config.DOMAIN)
}
