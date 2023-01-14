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
	

	//IF USER IS ALREADY SUBSCRIBED UNSUB THEM
	var alreadySubscribedCount int
	subscribedRowsError := conn.QueryRow(
		`select COUNT(*) from Subscribe where SubscriberID = ? and SubscribeeID= ?`, userId, subscribeeID.SubscribeeID).Scan(&alreadySubscribedCount)
	fmt.Println(subscribedRowsError)
	fmt.Println(alreadySubscribedCount)

	if (alreadySubscribedCount > 0) {
		//delete the subscription
		result, err := conn.Exec(
			`delete from Subscribe where SubscriberID = ? and SubscribeeID= ?`, userId, subscribeeID.SubscribeeID)
		fmt.Println(result)

		if err == nil {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Unsubscribed"})
			fmt.Println(err)
			return
	 	}
	} else {

	//IF USER IS NOT ALREADY SUBSCRIBED THEM
	result, err := conn.Exec(
		`insert into Subscribe (
			SubscriberID,
			SubscribeeID
        ) values (?, ?)`,
		userId,
		subscribeeID.SubscribeeID,
	)
	
	fmt.Println(result)
	fmt.Println(err)

	c.IndentedJSON(http.StatusOK, gin.H{"status" : 200, "message": "Subscribed successfully"})
	}
}