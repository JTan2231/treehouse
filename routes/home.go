package routes

import (
    "fmt"
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

func getSubscriptionArticles(userID int) ([]HomeArticle, error) {
    dbConn := db.GetDB()

    rows, err := dbConn.Query(
        `select
            a.Title,
            a.Slug,
            u.Username
        from Article a
        inner join Subscribe s on s.SubscriberID = ?
        inner join User u on u.UserID = s.SubscribeeID
        where a.UserID = s.SubscribeeID`, userID)

    if err != nil {
        return nil, err
    }

    var articles []HomeArticle

    if rows != nil {
        defer rows.Close()
        for rows.Next() {
            var article HomeArticle

            if err := rows.Scan(&article.Title, &article.Slug, &article.Username); err != nil {
                return nil, err
            }

            articles = append(articles, article)
        }
    }

    fmt.Println(articles);

    return articles, nil
}

func getExploreArticles(userID int) ([]HomeArticle, error) {
    dbConn := db.GetDB()

    rows, err := dbConn.Query(
        `select distinct
            a.Title,
            a.Slug,
            u.Username
        from Article a
        inner join User u on u.UserID = a.UserID and u.UserID != ?
        except
        select
            a.Title,
            a.Slug,
            u.Username
        from Article a
        inner join Subscribe s on s.SubscriberID = ?
        inner join User u on u.UserID = s.SubscribeeID
        where a.UserID = s.SubscribeeID`, userID, userID)

    if err != nil {
        return nil, err
    }

    var articles []HomeArticle

    if rows != nil {
        defer rows.Close()
        for rows.Next() {
            var article HomeArticle

            if err := rows.Scan(&article.Title, &article.Slug, &article.Username); err != nil {
                return nil, err
            }

            articles = append(articles, article)
        }
    }

    return articles, nil
}


func ServeHome(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")

	username, ok := session.Values["username"]
	if ok {
        userID := session.Values["userID"].(int)
        subscriptionArticles, err := getSubscriptionArticles(userID)
        exploreArticles, err := getExploreArticles(userID)

        if err != nil {
            c.IndentedJSON(400, gin.H{"errors": err})
            return
        }

		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"signedInUsername": username,
			"subscriptionArticles":         subscriptionArticles,
            "exploreArticles":          exploreArticles,
			"count":            len(subscriptionArticles),
		})
	} else {
		c.HTML(200, "404_redirect.tmpl", gin.H{
			"url": "/",
		})
	}
}
