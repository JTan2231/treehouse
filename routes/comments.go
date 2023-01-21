package routes

import (
	"encoding/json"
	"fmt"
    "time"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"treehouse/config"
	"treehouse/db"
	"treehouse/schema"
)

func CreateComment(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.IndentedJSON(http.StatusMethodNotAllowed, gin.H{})
		c.Abort()
		return
	}

	req, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		c.Abort()
	}

	newComment := schema.Comment{}

	json.Unmarshal(req, &newComment)

	newComment, err = addCommentToDB(newComment, c)

	session, _ := config.Store.Get(c.Request, "session")

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(400, gin.H{"message": err})
	} else {
		c.IndentedJSON(200, gin.H{
			"signedInUsername":   session.Values["username"],
			"comment_id": newComment.CommentID,
			"content":    newComment.Content,
		})
	}
}

// TODO
func verifyComment(comment schema.Comment) (schema.Comment, error) {
	return comment, nil
}

// TODO: move this to a separate file
// TODO: better error handling/DB constraints (duplicates, missing fields, etc.)
func addCommentToDB(comment schema.Comment, c *gin.Context) (schema.Comment, error) {
	conn := db.GetDB()
	newComment, err := verifyComment(comment)

    newComment.TimestampPosted = time.Now().Format("2006-01-02 15:04:05")

	if err != nil {
		return newComment, fmt.Errorf("CreateComment: %v", err)
	}

	// TODO: Check if slug exists in DB

	result, err := conn.Exec(
		`insert into Comment (
            UserID,
            ArticleID,
            ParentID,
            Content,
            TimestampPosted
        ) values (?, ?, ?, ?, ?)`,
		newComment.UserID,
		newComment.ArticleID,
		newComment.ParentID,
		newComment.Content,
        newComment.TimestampPosted,
	)
	if err != nil {
		return newComment, fmt.Errorf("CreateComment: %v", err)
	}

	commentId, err := result.LastInsertId()
	newComment.CommentID = int(commentId)

	if err != nil {
		return newComment, fmt.Errorf("CreateComment: %v", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return newComment, fmt.Errorf("CreateComment: %v", err)
	}

	return newComment, nil
}

type UserComment struct {
	Comment  schema.Comment
	Username string
}

type CommentTree struct {
	Comment  UserComment
	Children []*CommentTree
}

func treeFromComment(comment UserComment) CommentTree {
	return CommentTree{Comment: comment, Children: make([]*CommentTree, 0)}
}

func GetComments(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.IndentedJSON(http.StatusMethodNotAllowed, gin.H{})
		c.Abort()
		return
	}

	articleID, err := strconv.Atoi(c.Query("articleID"))
	if err != nil {
		c.IndentedJSON(400, gin.H{"message": "Bad request"})
		c.Abort()
		return
	}

	comments := queryComments(articleID)

	// construct n-ary comment tree from array of comments
	//
	// represent everything with array of CommentTrees,
	// use map of pointers to navigate
	// return the roots (CommentTrees without parents)

	trees := make([]CommentTree, 0)
	idPointerMap := make(map[int]*CommentTree)
	roots := make([]*CommentTree, 0)

	for i := 0; i < len(comments); i++ {
		trees = append(trees, treeFromComment(comments[i]))
		idPointerMap[comments[i].Comment.CommentID] = &trees[i]
	}

	for i := 0; i < len(comments); i++ {
		if current, ok := idPointerMap[comments[i].Comment.CommentID]; ok {
			if current.Comment.Comment.ParentID != nil {
				pid := *current.Comment.Comment.ParentID
				if parent, ok := idPointerMap[pid]; ok {
					parent.Children = append(parent.Children, current)
					idPointerMap[pid] = parent
				}
			} else {
				roots = append(roots, current)
			}
		}
	}

	c.IndentedJSON(200, gin.H{
		"comments": roots,
	})
}

func queryComments(articleID int) []UserComment {
	dbConn := db.GetDB()
	var comments []UserComment

	rows, err := dbConn.Query(
		`select
            c.CommentID,
            c.ArticleID,
            c.ParentID,
            c.UserID,
            c.Content,
			u.Username
        from Comment c 
		inner join User u 
		on c.UserID = u.UserID
		where c.articleID = ?`, articleID)

	if err != nil {
		fmt.Printf("error: %v", err)
		return comments
	}

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var comment UserComment

			if err := rows.Scan(
				&comment.Comment.CommentID,
				&comment.Comment.ArticleID,
				&comment.Comment.ParentID,
				&comment.Comment.UserID,
				&comment.Comment.Content,
				&comment.Username); err != nil {
				fmt.Printf("error: %v", err)
				return comments
			}

			comments = append(comments, comment)
		}
	}

	return comments
}
