package entity

import "time"

// 文章状态（生命周期）
const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
	ArticleStatusArchived  = "archived"
)

// 文章可见性（访问权限）
const (
	ArticleVisibilityPublic  = "public"  // 公开，所有人可见
	ArticleVisibilityPrivate = "private" // 私有，仅作者可见
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
	TagIDs          []int64 // 标签 ID 数组
	Status          string  // 状态：draft, published, archived
	Visibility      string  // 可见性：public, private
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

// ArticleLike 文章点赞。
type ArticleLike struct {
	ID          int64
	ArticleSlug string
	UserUUID    string
	CreatedAt   time.Time
}
