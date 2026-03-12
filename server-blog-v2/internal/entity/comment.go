package entity

import "time"

const (
	CommentStatusPending  = "pending"
	CommentStatusApproved = "approved"
	CommentStatusRejected = "rejected"
	CommentStatusSpam     = "spam"
)

// Comment 评论实体。
type Comment struct {
	ID          int64
	ArticleSlug string
	ParentID    *int64
	UserUUID    string
	Content     string
	Status      string
	Likes       int32
	IPAddress   *string
	UserAgent   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CommentLike 评论点赞。
type CommentLike struct {
	ID        int64
	CommentID int64
	UserUUID  string
	CreatedAt time.Time
}
