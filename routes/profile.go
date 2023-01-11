package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
    "treehouse/db"
    "treehouse/schema"
)

func ServeProfile(c *gin.Context) {
    var username = c.Param("username")
    dbConn := db.GetDB()

    rows, err := dbConn.Query(
        `select
            a.Title,
            a.Slug,
            u.UserID,
            u.Username
        from Article a
        inner join User u on u.Username = ?`, username)

    if err != nil {
        c.IndentedJSON(400, gin.H{ "errors" : "issue retrieving articles" })
    }

    defer rows.Close()

    var user schema.User
    var articles []schema.Article

    for rows.Next() {
        var article schema.Article

        if err := rows.Scan(&article.Title, &article.Slug, &user.UserID, &user.Username); err != nil {
            return
        }

        articles = append(articles, article)
    }

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"API_ROOT": config.API_ROOT,
        "articles": articles,
        "username": user.Username,
        "user_id": user.UserID,
	})
}
