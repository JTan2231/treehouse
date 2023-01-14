package routes

import (
	"encoding/json"
	"fmt"
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
			"username": session.Values["username"],
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

	// TODO: Check if slug exists in DB

	result, err := conn.Exec(
		`insert into Article (
            Title,
            Slug,
            Content,
            UserID
        ) values (?, ?, ?, ?)`,
		newArticle.Title,
		newArticle.Slug,
		newArticle.Content,
		newArticle.UserID,
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
	var username = c.Param("username")
	var slug = c.Param("slug")

	article := queryArticle(username, slug)

	c.HTML(http.StatusOK, "article_viewer.tmpl", gin.H{
		"title":   article.Title,
        "username": username,
		"content": strings.Split(article.Content, "\n"),
		"author":  username,
	})
}

func queryArticle(username string, slug string) schema.Article {
	conn := db.GetDB()

	var article schema.Article

	conn.QueryRow(`
            select
                Title,
                Content
            from Article a 
            inner join User u on u.Username = ? and u.UserID = a.UserID
            where a.Slug = ?
        `, username, slug).Scan(&article.Title, &article.Content)

	return article
}
