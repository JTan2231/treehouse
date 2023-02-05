package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"treehouse/config"
	"treehouse/db"
	"treehouse/schema"
)

func CreateNewUser(c *gin.Context) {
	req, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	newUser := schema.User{}
	json.Unmarshal(req, &newUser)

	//check if user is within db already
	check := checkIfUserExists(c, newUser)

	if check == 3 {
		c.IndentedJSON(400, gin.H{"status": 400, "message": "An account with that Username already exists"})
		return
	} else if check == 2 {
		c.IndentedJSON(400, gin.H{"status": 400, "message": "An account with that Email already exists"})
		return
	} else if check == 1 {
		c.IndentedJSON(400, gin.H{"status": 400, "message": "Account already exists with this email and username"})
		return
	}

	_, err = addUser(newUser)

	if err != nil {
		c.IndentedJSON(400, gin.H{"message": err})
		return
	} else {
		fmt.Println("User created successfully")
		//run a login auth here instead of routing them back to login,
		//this will allow them to be logged in after creating an account
		session, _ := config.Store.Get(c.Request, "session")
		conn := db.GetDB()
		stdmt := "SELECT UserID FROM User WHERE Username = ?"
		var userID int
		row := conn.QueryRow(stdmt, newUser.Username)
		row.Scan(&userID)

		session.Values["userID"] = userID
		session.Values["username"] = newUser.Username

		session.Save(c.Request, c.Writer)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func checkIfUserExists(c *gin.Context, newUser schema.User) int64 {
	conn := db.GetDB()

	//if it dosent exist return nil
	//3 username already exists
	//2 email already exists
	//1 account already exists
	//0 account does not exist

	var emailCount int
	email := conn.QueryRow(
		`select COUNT(Email) from User where Email = ?`, newUser.Email,
	)
	email.Scan(&emailCount)

	var usernameCount int
	username := conn.QueryRow(
		`select COUNT(Username) from User where Username = ?`, newUser.Username,
	)
	username.Scan(&usernameCount)

	fmt.Println(emailCount)
	fmt.Println(usernameCount)

	if emailCount != 0 && usernameCount != 0 {
		return 1
	} else if emailCount != 0 {
		return 2
	} else if usernameCount != 0 {
		return 3
	}

	return 0
}

func ServeNewUser(c *gin.Context) {
	c.HTML(http.StatusOK, "newuser.tmpl", gin.H{
		"API_ROOT": config.API_ROOT,
	})
}

// TODO
func verifyUser(user schema.User) (schema.User, error) {
	return user, nil
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

	profileUserId, _ := result.LastInsertId()
	defaultProfilePictureURL := "https://res.cloudinary.com/dubfvttoa/image/upload/v1675216906/defaultUser_fk4oe3.jpg"
	profileInsert, profileInsertErr := conn.Exec(
		`insert into Profile (
        UserID,
        ProfilePicture
    ) values (?, ?)`,
		profileUserId,
		defaultProfilePictureURL,
	)

	if profileInsertErr != nil {
		fmt.Println(profileInsert)
		fmt.Println(profileInsertErr)
		return 0, fmt.Errorf("addProfile: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}

	return id, nil
}
