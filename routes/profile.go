package routes

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strings"
	"time"
	"treehouse/config"
	"treehouse/db"
	"treehouse/schema"
)

type ProfileArticle struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	UserID          int64  `json:"userid"`
	Username        string `json:"username"`
	Subtitle        string `json:"subtitle"`
	TimestampPosted string `json:"timestampPosted"`
	Content         string `json:"content"`
	ReadTime        int    `json:"readTime"`
	ArticleID       int64  `json:"articleid"`
}

type Profile struct {
	Bio 	  string `json:"Bio"`
	TwitterURL string `json:"TwitterURL"`
}

func GetEditProfile(c *gin.Context) {
	db := db.GetDB()
	session, _ := config.Store.Get(c.Request, "session")
	localusername := session.Values["username"]
	localUserID := session.Values["userID"]

	var profile Profile
	var profilePicURL string
	db.QueryRow(`SELECT Bio, TwitterURL, ProfilePicture FROM Profile WHERE UserID = ?`, localUserID).Scan(&profile.Bio, &profile.TwitterURL, &profilePicURL)

	c.HTML(http.StatusOK, "editProfile.tmpl", gin.H{
		"username": localusername,
		"twitterURL": profile.TwitterURL,
		"profilePicture": profilePicURL,
		"bio": profile.Bio,
	})
}

func EditProfile(c *gin.Context)  {
	dbConn := db.GetDB()
	session, _ := config.Store.Get(c.Request, "session")
	userid := session.Values["userID"]

	var profile Profile
	c.BindJSON(&profile)

	result, err := dbConn.Exec(`UPDATE Profile SET Bio = ?, TwitterURL = ? WHERE UserID = ?`, profile.Bio, profile.TwitterURL, userid)

	
	if(err != nil){
		fmt.Println(result)
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Profile Update Failed",
		})
	}

	//rehydrate html page with all new data from db
	c.IndentedJSON(http.StatusOK, gin.H{
		"bio" : profile.Bio,
	})
}

func GetHeaderProfilePic(c *gin.Context) {
	dbConn := db.GetDB()
	session, _ := config.Store.Get(c.Request, "session")

	//var username = session.Values["username"]
	var localuserID = session.Values["userID"]


	var profilePicURL string
	_ = dbConn.QueryRow(`select ProfilePicture from Profile where UserID = ?`, localuserID).Scan(&profilePicURL)
	fmt.Println(profilePicURL)

	c.JSON(http.StatusOK, gin.H{
		"profilePicURL": profilePicURL,
	})
}

func GetLocalUserName(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	localusername := session.Values["username"]
	c.JSON(http.StatusOK, gin.H{
		"username": localusername,
	})
}

func ServeProfile(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	localusername := session.Values["username"]
	localuserID := session.Values["userID"]

	var username = c.Param("username")
	dbConn := db.GetDB()

	var profileUserID = 0
	if err := dbConn.QueryRow("select UserID from User where Username = ?", username).Scan(&profileUserID); err != nil {
		if err == sql.ErrNoRows {
			c.HTML(404, "404_redirect.tmpl", gin.H{
				"url": "/home",
			})
			return
		}
	}

	rows, err := dbConn.Query(
		`select
            a.Title,
            a.Slug,
			a.Subtitle,
			a.TimestampPosted,
			a.Content,
            u.Username,
			a.ArticleID
        from Article a
        inner join User u on u.UserID = a.UserID
        where u.Username = ?`, username)

	if err != nil {
		c.IndentedJSON(400, gin.H{"errors": err})
		return
	}

	//diff select statement for user id if  it is getting passed incorrectly
	//bug if user does not have an article
	var user schema.User
	var articles []ProfileArticle

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var article ProfileArticle

			if err := rows.Scan(&article.Title, &article.Slug, &article.Subtitle, &article.TimestampPosted, &article.Content, &article.Username, &article.ArticleID); err != nil {
				return
			}

			//formatting date to english
			engDate := (strings.Split(article.TimestampPosted, " ")[0])
			t, _ := time.Parse("2006-01-02", engDate)
			article.TimestampPosted = t.Format("January 02, 2006")

			//calculating read time
			var words int
			var length int
			length = len(strings.Split(article.Content, " "))

			words = int(math.Ceil(float64(length) / 238))
			article.ReadTime = words

			articles = append(articles, article)
		}
	}

	//getting user id without being dependent on if they have an aritcle or not
	userIDAndNameRow := dbConn.QueryRow("select UserID,Username from User where Username = ?", username).Scan(&user.UserID, &user.Username)

	if userIDAndNameRow != nil {
		fmt.Println(userIDAndNameRow)
	}

	rows, err = dbConn.Query(
		`select
            u.Username
        from User u
        inner join Subscribe s on s.SubscriberID = ? and u.UserID = s.SubscribeeID`, profileUserID)

	if err != nil {
		c.IndentedJSON(400, gin.H{"errors": err})
		return
	}

	var subscriptions []string

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var subscription string

			if err := rows.Scan(&subscription); err != nil {
				return
			}

			subscriptions = append(subscriptions, subscription)
		}
	}

	rows, err = dbConn.Query(
		`select
            a.Title,
            a.Slug,
            u.UserID,
            u.Username
        from Article a
        inner join User u on u.UserID = a.UserID
		inner join Favorite f on f.UserID = ? and a.ArticleID = f.ArticleID`, profileUserID)

	if err != nil {
		c.IndentedJSON(400, gin.H{"errors": err})
		return
	}

	var favorites []ProfileArticle

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var favorite ProfileArticle

			if err := rows.Scan(&favorite.Title, &favorite.Slug, &favorite.UserID, &favorite.Username); err != nil {
				return
			}

			favorites = append(favorites, favorite)
		}
	}

	check := (localusername == username)

	alreadySubscribedBool := false
	var alreadySubscribedCount int

	// TODO: the following query can probably be removed/moved/merged
	// checking if they are already subscribed, if so, set alreadySubscribed to true
	subscribedRowsError := dbConn.QueryRow(
		`select COUNT(*) from Subscribe where SubscriberID = ? and SubscribeeID = ?`, localuserID, profileUserID).Scan(&alreadySubscribedCount)

	if subscribedRowsError != nil {
		fmt.Println(subscribedRowsError)
	}

	alreadySubscribedBool = alreadySubscribedCount > 0

	var profilePicURL string
	_ = dbConn.QueryRow(`select ProfilePicture from Profile where UserID = ?`, profileUserID).Scan(&profilePicURL)
	fmt.Println(profilePicURL)


	//JOEY PEEP THIS THERE IS A PROBLEM WITH SCANNING SINCE I ADDED NULL TO THE COLUMN IN TEH SQL SCRIPT
	var bio sql.NullString
	result := dbConn.QueryRow(`select Bio from Profile where UserID = ?`, profileUserID).Scan(&bio)
	fmt.Println(bio)
	fmt.Println(result)


	var twitterURL string
	_ = dbConn.QueryRow(`select TwitterURL from Profile where UserID = ?`, profileUserID).Scan(&twitterURL)
	fmt.Println(twitterURL)


	var twitterCheck bool
	if(twitterURL == "") {
		twitterCheck = false
	} else {
		twitterCheck = true
	}

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"API_ROOT":          config.API_ROOT,
		"articles":          articles,
		"subscriptions":     subscriptions,
		"favorites":         favorites,
		"signedInUsername":  localusername,
		"profileUsername":   username,
		"user_id":           profileUserID,
		"check":             check,
		"alreadySubscribed": alreadySubscribedBool,
		"profilePicURL" : profilePicURL,
		"bio" : bio.String,
		"twitterURL" : twitterURL,
		"twitterCheck" : twitterCheck,
	})
}
