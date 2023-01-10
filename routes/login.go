package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
    "treehouse/db"
    "golang.org/x/crypto/bcrypt"
    config "treehouse/config"
)

type LoginUser struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func AuthenticateLogin(c *gin.Context) {
    conn := db.GetDB()
    user := LoginUser{}
    c.BindJSON(&user)


    var userID int
    var hash string

    stdmt := "SELECT UserID, Password FROM User WHERE Username = ?"
    row := conn.QueryRow(stdmt, user.Username)
    err := row.Scan(&userID, &hash)
    if err != nil {
        fmt.Println(err)
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
    //check db to see if the user is within it
    //if they are create a session and log them in
    //query their username, find their userID
    //compare theire hash to the password inputted


    //need middleware for routes that have /users which whill check if a user is authenticated
}

func ServeLogin(c *gin.Context) { 
    c.HTML(http.StatusOK, "login.tmpl", gin.H{
        "API_ROOT": config.API_ROOT,
    })
}
