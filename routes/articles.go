package routes

import (
    "io/ioutil"
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gin-gonic/gin"

    "treehouse/schema"
    "treehouse/db"
)

func CreateArticle(c *gin.Context) {
    req, err := ioutil.ReadAll(c.Request.Body)

    //print the json request sent
    fmt.Println(string(req))

    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Bad request" })
        return
    }

    newArticle := schema.Article{}
    json.Unmarshal(req, &newArticle)

    _, err = createArticle(newArticle)

    if err != nil {
        fmt.Println(err)
        c.IndentedJSON(400, gin.H{ "message" : err })
    } else {
        c.IndentedJSON(http.StatusOK, gin.H{ "message": "Success" })
    }
}

// TODO
func verifyArticle(article  schema.Article) (schema.Article, error) {
    return article, nil
}

// TODO: better error handling/DB constraints (duplicates, missing fields, etc.)
func createArticle(article schema.Article) (int64, error) {
    conn := db.GetDB()

    newArticle, err := verifyArticle(article)
    newArticle.UserID = 1

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
        "title": article.Title,
        "content": article.Content,
    })
}

func queryArticle(username string, title string) (schema.Article) {
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
