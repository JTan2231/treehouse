package routes

import (
    "io/ioutil"
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gin-gonic/gin"

    schema "treehouse/schema"
)

func CreateNewUser(c *gin.Context) {
    req, err := ioutil.ReadAll(c.Request.Body)
 
    //print the json request sent   
    fmt.Println(string(req))
    
    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Bad request" })
        return
    }
        
    newUser := schema.User{}
    json.Unmarshal(req, &newUser)
     
    //addUser(newUser) to DB and hash pw with bcrypt
    
    c.IndentedJSON(http.StatusOK, gin.H{ "message": "User created successfully" })
}

func ServeNewUser(c *gin.Context) {
    c.HTML(http.StatusOK, "newuser.tmpl", gin.H{})
}
