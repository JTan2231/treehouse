package routes

import (
	"github.com/gin-gonic/gin"
	"treehouse/config"
    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"os"
	"strconv"
	"treehouse/db"
)


func UploadProfilePic(c *gin.Context) {
	file, _, err := c.Request.FormFile("profilePicture")
	cld , _ := cloudinary.NewFromParams("dubfvttoa", os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	session, _ := config.Store.Get(c.Request, "session")
	userID := session.Values["userID"].(int)
	userIDString := strconv.Itoa(userID)

	// Upload the file to cloudinary
	uploadResult, err := cld.Upload.Upload(c, file, uploader.UploadParams{PublicID: userIDString})

	if err != nil {
		c.IndentedJSON(500, gin.H{})
		return
	}

	// Get the profile picture URL
	profilePicURL := uploadResult.SecureURL
	
	// Update the profile picture URL in the database
	dbConn := db.GetDB()
	_, err = dbConn.Exec(`UPDATE Profile SET ProfilePicture = ? WHERE UserID = ?`, profilePicURL, userID)


	c.IndentedJSON(200, gin.H{
		"username" : session.Values["username"],
		"message" : "success",
		"profilePicURL": profilePicURL,
	})
}