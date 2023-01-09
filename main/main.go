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

    // TODO: this is gross. gotta be a better way to organize routes/endpoints than below
    router.GET("/users/:username/:title", routes.GetArticle)
    router.GET("/users/:username/create-article", routes.GetCreateArticle)

    router.POST("/articles", routes.CreateArticle)
    router.POST("/newuser", routes.CreateNewUser)

    router.Run(config.DOMAIN)
}
