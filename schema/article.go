package schema

type Article struct {
	ArticleID int    `json:"articleid"`
	UserID    int    `json:"userid"`
	Title     string `json:"title"`
	Slug     string `json:"slug"`
	Content   string `json:"content"`
}
