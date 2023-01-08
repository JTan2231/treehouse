package routes

import (
    "io/ioutil"
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "treehouse/config"
    "treehouse/db"
    "treehouse/schema"
)

func CreateNewUser(c *gin.Context) {
    req, err := ioutil.ReadAll(c.Request.Body)

    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Bad request" })
        return
    }
    
    newUser := schema.User{}
    json.Unmarshal(req, &newUser)

    _, err = addUser(newUser)

    if err != nil {
        c.IndentedJSON(400, gin.H{ "message": err })
    } else {
        fmt.Println("User created successfully")
        c.IndentedJSON(http.StatusOK, gin.H{ "message": "User created successfully" })
    }
}

func ServeNewUser(c *gin.Context) {
    c.HTML(http.StatusOK, "newuser.tmpl", gin.H{
        "API_ROOT": config.API_ROOT,
    })
}

// TODO
func verifyUser(user schema.User) (schema.User, error) {
    return user, nil;
}

func addUser(user schema.User) (int64, error) {
    conn := db.GetDB()

    newUser, err := verifyUser(user)

    if err != nil {
        return 0, fmt.Errorf("addUser: %v", err)
    }

    pass := []byte(newUser.Password)
    hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

    result, err := conn.Exec(
        `insert into User (
            Username,
            Email,
            Password
        ) values (?, ?, ?)`,
        newUser.Username,
        newUser.Email,
        hashed,
    )

    if err != nil {
        return 0, fmt.Errorf("addUser: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("addUser: %v", err)
    }

    return id, nil
}
