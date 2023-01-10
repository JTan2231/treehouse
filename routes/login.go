package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
    "treehouse/db"
    "treehouse/schema"
    "golang.org/x/crypto/bcrypt"
    config "treehouse/config"
)

func AuthenticateLogin(c *gin.Context) {
    conn := db.GetDB()
    user := schema.LoginUser{}
    c.BindJSON(&user)


    var userID int
    var hash string

    stdmt := "SELECT UserID, Password FROM User WHERE Username = ?"
    row := conn.QueryRow(stdmt, user.Username)
    err := row.Scan(&userID, &hash)
    if err != nil {
        fmt.Println(err)
        fmt.Println("user not found")
        c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Bad request" })
        return;
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password))
        if(err == nil) {
            fmt.Println("passwords match")
            session, _ := config.Store.Get(c.Request, "session")
            session.Values["userID"] = userID
            session.Save(c.Request, c.Writer)
            c.IndentedJSON(http.StatusOK, gin.H{ "message": "it works" })
            return
        } else {
            fmt.Println("passwords do not match")
            c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Bad request" })
            return
        }
}

func AuthRequired(c *gin.Context) {
    session, _ := config.Store.Get(c.Request, "session")
    _,ok := session.Values["userID"]
    if !ok {
        c.HTML(http.StatusForbidden, "login.tmpl", nil)
        c.Abort()
        return
    }
    c.Next()
}

func ServeLogin(c *gin.Context) { 
    c.HTML(http.StatusOK, "login.tmpl", gin.H{
        "API_ROOT": config.API_ROOT,
    })
}
