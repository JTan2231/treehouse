package routes

import (
	"encoding/json"
	"fmt"
    "time"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"treehouse/config"
	"treehouse/db"
	"treehouse/schema"
)

func CreateArticle(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.IndentedJSON(http.StatusMethodNotAllowed, gin.H{})
		c.Abort()
		return
	}

	req, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		c.Abort()
	}

	newArticle := schema.Article{}

	json.Unmarshal(req, &newArticle)

	newArticle, err = addArticleToDB(newArticle, c)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(400, gin.H{"message": err})
	} else {
		session, _ := config.Store.Get(c.Request, "session")

		c.IndentedJSON(200, gin.H{
			"slug":     newArticle.Slug,
			"signedInUsername": session.Values["username"],
		})
	}
}

func GetCreateArticle(c *gin.Context) {
	c.HTML(http.StatusOK, "create_article.tmpl", gin.H{
		"API_ROOT": config.API_ROOT,
	})
}

// TODO
func verifyArticle(article schema.Article) (schema.Article, error) {
	return article, nil
}

// TODO: move this to a separate file
func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

// TODO: better error handling/DB constraints (duplicates, missing fields, etc.)
func addArticleToDB(article schema.Article, c *gin.Context) (schema.Article, error) {
	conn := db.GetDB()
	newArticle, err := verifyArticle(article)

	session, _ := config.Store.Get(c.Request, "session")

	//error checking to make sure this value is not null******
	idOfUser, ok := session.Values["userID"]

	if !ok {
		fmt.Printf("you are not logged in")
		return newArticle, fmt.Errorf("createArticle: %v", err)
	}

	if err != nil {
		return newArticle, fmt.Errorf("createArticle: %v", err)
	}

	newArticle.UserID = idOfUser.(int)
	newArticle.Slug = strip(newArticle.Title)
	newArticle.Slug = strings.ToLower(strings.ReplaceAll(newArticle.Slug, " ", "-"))

    newArticle.TimestampPosted = time.Now().Format("2006-01-02 15:04:05")
    fmt.Println("New article at " + newArticle.TimestampPosted + " UTC!")

	// TODO: Check if slug exists in DB

	result, err := conn.Exec(
		`insert into Article (
            Title,
            Subtitle,
            Slug,
            Content,
            UserID,
            TimestampPosted
        ) values (?, ?, ?, ?, ?, ?)`,
		newArticle.Title,
        newArticle.Subtitle,
		newArticle.Slug,
		newArticle.Content,
		newArticle.UserID,
        newArticle.TimestampPosted,
	)

	if err != nil {
		return newArticle, fmt.Errorf("createArticle: %v", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return newArticle, fmt.Errorf("createArticle: %v", err)
	}

	return newArticle, nil
}

func GetArticle(c *gin.Context) {
	var authorUsername = c.Param("username")
	var slug = c.Param("slug")
	session, _ := config.Store.Get(c.Request, "session")
	dbConn := db.GetDB()

	article := queryArticle(authorUsername, slug)

	alreadyFavoritedBool := false
	var alreadyFavoritedCount int

	favoriteRowsError := dbConn.QueryRow(
		`select COUNT(*) from Favorite where UserID = ? and ArticleID= ?`, session.Values["userID"], article.ArticleID).Scan(&alreadyFavoritedCount)

	if favoriteRowsError != nil {
		fmt.Println(favoriteRowsError)
	}

    fmt.Println("ARTICLE ID: ", article.ArticleID)

	alreadyFavoritedBool = alreadyFavoritedCount > 0

	c.HTML(http.StatusOK, "article_viewer.tmpl", gin.H{
		"content":          strings.Split(article.Content, "\n"),
		"localUserID":      session.Values["userID"],
		"alreadyFavorited": alreadyFavoritedBool,
		"articleID":        article.ArticleID,
		"title":            article.Title,
        "subtitle": article.Subtitle,
        "timestamp": article.TimestampPosted,
		"authorUsername":   authorUsername,
		"signedInUsername":    session.Values["username"],
	})
}

func queryArticle(username string, slug string) schema.Article {
	conn := db.GetDB()

	var article schema.Article

    err := conn.QueryRow(`
            select
                ArticleID,
                Title,
                Subtitle,
                Content,
                TimestampPosted
            from Article a 
            inner join User u on u.Username = ? and u.UserID = a.UserID
            where a.Slug = ?
        `, username, slug).Scan(&article.ArticleID, &article.Title, &article.Subtitle, &article.Content, &article.TimestampPosted)

    if err != nil {
        fmt.Println("queryArticle: ", err)
    }

	return article
}
