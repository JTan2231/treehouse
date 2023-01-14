package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
	"treehouse/db"
	"treehouse/schema"
)

func ServeProfile(c *gin.Context) {
	var username = c.Param("username")
	dbConn := db.GetDB()

	var i = 0
	if err := dbConn.QueryRow("select 1 from User where Username = ?", username).Scan(&i); err != nil {
		if err == sql.ErrNoRows {
			c.HTML(404, "404_redirect.tmpl", gin.H{
				"url": "/home",
			})
			return
		}
	}

	rows, err := dbConn.Query(
		`select
            a.Title,
            a.Slug,
            u.UserID,
            u.Username
        from Article a
        inner join User u on u.UserID = a.UserID
        where u.Username = ?`, username)

	if err != nil {
		c.IndentedJSON(400, gin.H{"errors": err})
		return
	}

	var user schema.User
	var articles []schema.Article

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var article schema.Article

			if err := rows.Scan(&article.Title, &article.Slug, &user.UserID, &user.Username); err != nil {
				return
			}

			articles = append(articles, article)
		}
	}

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"API_ROOT": config.API_ROOT,
		"articles": articles,
		"username": username,
		"user_id":  user.UserID,
	})
}
