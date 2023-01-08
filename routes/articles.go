package routes

import (
    "io/ioutil"
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
)

func CreateArticle(c *gin.Context) {
    req, err := ioutil.ReadAll(c.Request.Body)

    //print the json request sent
    fmt.Println(string(req))

    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Bad request" })
        return
    }

    //newUser := User{}
    //json.Unmarshal(req, &newUser)

    //addUser(newUser) to DB and hash pw with bcrypt

    c.IndentedJSON(http.StatusOK, gin.H{ "message": "Success" })
    c.HTML(http.StatusOK, "newuser.tmpl", gin.H{
        "post"  : "post test",
    })
}

func serveNewUser(c *gin.Context) {
    c.HTML(http.StatusOK, "newuser.tmpl", gin.H{})
}
