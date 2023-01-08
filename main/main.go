package main

import (
    "github.com/gin-gonic/gin"

    routes "treehouse/routes"
    config "treehouse/config"
)

const (
    ContentTypeBinary = "application/octet-stream"
    ContentTypeForm   = "application/x-www-form-urlencoded"
    ContentTypeJSON   = "application/json"
    ContentTypeHTML   = "text/html; charset=utf-8"
    ContentTypeText   = "text/plain; charset=utf-8"
)

func main() {
    //initDB()

    router := gin.Default()
    router.LoadHTMLGlob("templates/*")
    router.Static("/styles", "./styles")

    router.GET("/", routes.ServeLogin) // TEMP: get an actual homepage later
    router.POST("/articles", routes.CreateArticle)

    router.GET("/makePost", routes.MakePost)

    router.Run(config.DOMAIN)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
/*
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")

    // Loop through the list of albums, looking for
    // an album whose ID value matches the parameter.
    for _, a := range albums {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}*/
