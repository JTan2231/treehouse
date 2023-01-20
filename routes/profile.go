package routes

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
	"treehouse/db"
	"treehouse/schema"
)

type ProfileArticle struct {
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	UserID   int64  `json:"userid"`
	Username string `json:"username"`
}

func ServeProfile(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	localusername := session.Values["username"]
	localuserID := session.Values["userID"]

	var username = c.Param("username")
	dbConn := db.GetDB()

	var profileUserID = 0
	if err := dbConn.QueryRow("select UserID from User where Username = ?", username).Scan(&profileUserID); err != nil {
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
            u.Username
        from Article a
        inner join User u on u.UserID = a.UserID
        where u.Username = ?`, username)

	if err != nil {
		c.IndentedJSON(400, gin.H{"errors": err})
		return
	}

	//diff select statement for user id if  it is getting passed incorrectly
	//bug if user does not have an article
	var user schema.User
	var articles []ProfileArticle

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var article ProfileArticle

			if err := rows.Scan(&article.Title, &article.Slug, &article.UserID, &article.Username); err != nil {
				return
			}

			articles = append(articles, article)
		}
	}

	//getting user id without being dependent on if they have an aritcle or not
	userIDAndNameRow := dbConn.QueryRow("select UserID,Username from User where Username = ?", username).Scan(&user.UserID, &user.Username)

	if userIDAndNameRow != nil {
		fmt.Println(userIDAndNameRow)
	}

	rows, err = dbConn.Query(
		`select
            u.Username
        from User u
        inner join Subscribe s on s.SubscriberID = ? and u.UserID = s.SubscribeeID`, profileUserID)

	if err != nil {
		c.IndentedJSON(400, gin.H{"errors": err})
		return
	}

	var subscriptions []string

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var subscription string

			if err := rows.Scan(&subscription); err != nil {
				return
			}

			subscriptions = append(subscriptions, subscription)
		}
	}

	check := (localusername == username)

	alreadySubscribedBool := false
	var alreadySubscribedCount int

	// TODO: the following query can probably be removed/moved/merged
	// checking if they are already subscribed, if so, set alreadySubscribed to true
	subscribedRowsError := dbConn.QueryRow(
		`select COUNT(*) from Subscribe where SubscriberID = ? and SubscribeeID = ?`, localuserID, profileUserID).Scan(&alreadySubscribedCount)

	if subscribedRowsError != nil {
		fmt.Println(subscribedRowsError)
	}

	alreadySubscribedBool = alreadySubscribedCount > 0

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"API_ROOT":          config.API_ROOT,
		"articles":          articles,
		"subscriptions":     subscriptions,
		"username":          localusername,
		"profileUsername":   username,
		"user_id":           profileUserID,
		"check":             check,
		"alreadySubscribed": alreadySubscribedBool,
	})
}
