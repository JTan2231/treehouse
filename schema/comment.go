package schema

type Comment struct {
	CommentID int    `json:"commentid"`
	UserID    int    `json:"userid"`
	ArticleID int    `json:"articleid"`
	ParentID  *int   `json:"parentid"`
	Content   string `json:"content"`
}
