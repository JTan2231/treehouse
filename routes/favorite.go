package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treehouse/db"
)

type Favorite struct {
	UserID    int `json:"userID"`
	ArticleID int `json:"articleID"`
}

func FavoriteArticle(c *gin.Context) {
	db := db.GetDB()

	favorite := Favorite{}
	c.BindJSON(&favorite)

	var alreadyFavoriteCount int
	_ = db.QueryRow(
		`select COUNT(*) from Favorite where UserID = ? and ArticleID= ?`, favorite.UserID, favorite.ArticleID).Scan(&alreadyFavoriteCount)


	if alreadyFavoriteCount > 0 {
		_, err := db.Exec(
			`delete from Favorite where UserID = ? and ArticleID= ?`, favorite.UserID, favorite.ArticleID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error unfavoriting article"})
			return
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"status": 200, "message": "Successfully unfavorited article"})
		}

	} else {
		_, err := db.Exec("INSERT INTO Favorite (ArticleID, UserID) VALUES (?, ?)", favorite.ArticleID, favorite.UserID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error favoriting article"})
			return
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"status": 200, "message": "Successfully favorited article"})
			return
		}
	}
}
