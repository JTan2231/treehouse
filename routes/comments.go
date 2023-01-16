package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"treehouse/db"
	"treehouse/config"
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

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(400, gin.H{"message": err})
	} else {
		c.IndentedJSON(200, gin.H{
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

	if err != nil {
		return newComment, fmt.Errorf("CreateComment: %v", err)
	}

	// TODO: Check if slug exists in DB

	result, err := conn.Exec(
		`insert into Comment (
            UserID,
            ArticleID,
            ParentID,
            Content
        ) values (?, ?, ?, ?)`,
		newComment.UserID,
		newComment.ArticleID,
		newComment.ParentID,
		newComment.Content,
	)

	if err != nil {
		return newComment, fmt.Errorf("CreateComment: %v", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return newComment, fmt.Errorf("CreateComment: %v", err)
	}

	return newComment, nil
}

type CommentTree struct {
	Comment  schema.Comment
	Children []*CommentTree
}

func treeFromComment(comment schema.Comment) CommentTree {
	return CommentTree{Comment: comment, Children: make([]*CommentTree, 0)}
}

func GetComments(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.IndentedJSON(http.StatusMethodNotAllowed, gin.H{})
		c.Abort()
		return
	}

	session, _ := config.Store.Get(c.Request, "session")

	userID := session.Values["userID"].(int)
	comments := queryComments(userID)

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
		idPointerMap[comments[i].CommentID] = &trees[i]
	}

	for i := 0; i < len(comments); i++ {
		if current, ok := idPointerMap[comments[i].CommentID]; ok {
			if current.Comment.ParentID != nil {
				pid := *current.Comment.ParentID
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

func queryComments(userID int) []schema.Comment {
	dbConn := db.GetDB()
	var comments []schema.Comment

	rows, err := dbConn.Query(
		`select
            CommentID,
            ArticleID,
            ParentID,
            UserID,
            Content
        from Comment c where c.UserID = ?`, userID)

	if err != nil {
		fmt.Printf("error: %v", err)
		return comments
	}

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var comment schema.Comment

			if err := rows.Scan(
				&comment.CommentID,
				&comment.ArticleID,
				&comment.ParentID,
				&comment.UserID,
				&comment.Content); err != nil {
				fmt.Printf("error: %v", err)
				return comments
			}

			comments = append(comments, comment)
		}
	}

	return comments
}
