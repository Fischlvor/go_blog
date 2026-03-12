package entity

import "time"

const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
	ArticleStatusArchived  = "archived"
)

// Article 文章实体。
type Article struct {
	ID              int64
	Title           string
	Slug            string
	Excerpt         *string
	Content         string
	FeaturedImage   *string
	AuthorUUID      string // 作者 UUID
	CategoryID      int64
	Status          string
	ReadTime        *string
	Views           int32
	Likes           int32
	IsFeatured      bool
	MetaTitle       *string
	MetaDescription *string
	PublishedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ArticleTag 文章标签关联。
type ArticleTag struct {
	ArticleSlug string
	TagID       int64
}

// ArticleLike 文章点赞。
type ArticleLike struct {
	ID          int64
	ArticleSlug string
	UserUUID    string
	CreatedAt   time.Time
}
