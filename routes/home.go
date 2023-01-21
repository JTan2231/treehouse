package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
	"treehouse/db"
)

type HomeArticle struct {
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Username string `json:"username"`
}

func ServeHome(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")

	username, ok := session.Values["username"]
	if ok {
		dbConn := db.GetDB()

		userID := session.Values["userID"].(int)
		rows, err := dbConn.Query(
			`select
                a.Title,
                a.Slug,
                u.Username
            from Article a
            inner join Subscribe s on s.SubscriberID = ?
            inner join User u on u.UserID = s.SubscribeeID
            where a.UserID = s.SubscribeeID`, userID)
		//temp table with all articles from users you are subscribed to
		//understand

		if err != nil {
			c.IndentedJSON(400, gin.H{"errors": err})
			return
		}

		var articles []HomeArticle

		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var article HomeArticle

				if err := rows.Scan(&article.Title, &article.Slug, &article.Username); err != nil {
					return
				}

				articles = append(articles, article)
			}
		}

		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"signedInUsername": username,
			"articles": articles,
            "count": len(articles),
		})
	} else {
		c.HTML(200, "404_redirect.tmpl", gin.H{
			"url": "/",
		})
	}
}
