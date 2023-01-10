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

    //router.GET("/users/:username", routes.GetProfile)
    router.GET("/", routes.ServeLogin) // TODO: get an actual homepage
    router.GET("/users/:username/:title", routes.GetArticle)
    router.GET("/users/:username/createarticle", routes.GetCreateArticle)
    router.GET("/newuser" ,routes.ServeNewUser)
    
    router.POST("/login", routes.AuthenticateLogin)
    router.POST("/newuser", routes.CreateNewUser)
    router.POST("/articles", routes.CreateArticle)


    //authRouter := router.Group("/user",auth)

    

    router.Run(config.DOMAIN)
}
