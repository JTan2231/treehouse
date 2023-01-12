package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/config"
	"treehouse/db"
	"fmt"
)

//user sends a post request to subscribe
//with the subscribeeID(which is already passed into frontend);
//the local user subscribes himself to the profile he is on

type SubscribeID struct {
	SubscribeeID int `json:"subscribeeID"`
}

func SubscribeToUser(c *gin.Context) {
	conn := db.GetDB()
	session, _ := config.Store.Get(c.Request, "session")
	userId := session.Values["userID"]


	subscribeeID := SubscribeID{}
	c.BindJSON(&subscribeeID)

	//subscriber ID from cookie
	//subscribeeID from request body

	result, err := conn.Exec(
		`insert into Subscribe (
			SubscriberID,
			SubscribeeID
        ) values (?, ?)`,
		userId,
		subscribeeID.SubscribeeID,
	)
	
	fmt.Println(result)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		fmt.Println(err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status" : 200, "message": "user subscribed successfully"})
}