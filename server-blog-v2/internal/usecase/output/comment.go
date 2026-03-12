package output

import "time"

// Comment 评论。
type Comment struct {
	ID          int64       `json:"id"`
	ArticleSlug string      `json:"article_slug"`
	Content     string      `json:"content"`
	User        CommentUser `json:"user"`
	ParentID    *int64      `json:"parent_id"`
	Children    []Comment   `json:"children"`
	CreatedAt   time.Time   `json:"created_at"`
}

// CommentUser 评论用户信息。
type CommentUser struct {
	UUID     string `json:"uuid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// UserBrief 用户简要信息。
type UserBrief struct {
	UUID     string `json:"uuid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CommentAdmin 评论（管理端）。
type CommentAdmin struct {
	ID           int64       `json:"id"`
	ArticleSlug  string      `json:"article_slug"`
	ArticleTitle string      `json:"article_title"`
	Content      string      `json:"content"`
	User         CommentUser `json:"user"`
	ParentID     *int64      `json:"parent_id"`
	CreatedAt    time.Time   `json:"created_at"`
}
