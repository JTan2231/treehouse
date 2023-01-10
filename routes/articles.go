package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
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

	_, err = addArticleToDB(newArticle, c)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(400, gin.H{"message": err})
	} else {
		//send a 200 response and check for 200 on frontened then redirect
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

// TODO: better error handling/DB constraints (duplicates, missing fields, etc.)
func addArticleToDB(article schema.Article, c *gin.Context) (int64, error) {
	conn := db.GetDB()
	newArticle, err := verifyArticle(article)

	session, _ := config.Store.Get(c.Request, "session")

	//error checking to make sure this value is not null******
	idOfUser, ok := session.Values["userID"]
	newArticle.UserID = idOfUser.(int)

	if !ok {
		fmt.Printf("you are not logged in")
		return -1, fmt.Errorf("createArticle: %v", err)
	}

	if err != nil {
		return 0, fmt.Errorf("createArticle: %v", err)
	}

	result, err := conn.Exec(
		`insert into Article (
            Title,
            Content,
            UserID
        ) values (?, ?, ?)`,
		newArticle.Title,
		newArticle.Content,
		newArticle.UserID,
	)

	if err != nil {
		return 0, fmt.Errorf("createArticle: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createArticle: %v", err)
	}

	fmt.Printf("NEW ARTICLE ADDED ID: %v", newArticle.Content)

	return id, nil
}

func GetArticle(c *gin.Context) {
	var username = c.Param("username")
	var title = c.Param("title")

	article := queryArticle(username, title)

	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":   article.Title,
		"content": article.Content,
	})
}

func queryArticle(username string, title string) schema.Article {
	conn := db.GetDB()

	var article schema.Article

	conn.QueryRow(`
            select
                Title,
                Content
            from Article a 
            inner join User u on u.Username = ? and u.UserID = a.UserID
            where a.Title = ?
        `, username, title).Scan(&article.Title, &article.Content)

	fmt.Printf("\n\n\n%v\n\n%v\n\n\n", article.Title, article.Content)

	return article
}
