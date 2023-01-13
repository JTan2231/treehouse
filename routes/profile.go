package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
	"treehouse/db"
	"treehouse/schema"
	"fmt"
)

func ServeProfile(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	localusername := session.Values["username"]
	localuserID := session.Values["userID"]

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
		c.IndentedJSON(400, gin.H{"errors": "issue retrieving articles"})
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

	check := (localusername == user.Username)

	alreadySubscribedBool := false
	var alreadySubscribedCount int


	//checking if they are already subscribed, if so, set alreadySubscribed to true
	subscribedRowsError := dbConn.QueryRow(
		`select COUNT(*) from Subscribe where SubscriberID = ? and SubscribeeID= ?`, localuserID, user.UserID).Scan(&alreadySubscribedCount)

	if subscribedRowsError != nil {
		fmt.Println(subscribedRowsError)
	}

	if (alreadySubscribedCount > 0){
		alreadySubscribedBool = true
	}

	

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"API_ROOT": config.API_ROOT,
		"articles": articles,
		"username": user.Username,
		"user_id":  user.UserID,
		"check": check,
		"alreadySubscribed": alreadySubscribedBool,
	})
}
