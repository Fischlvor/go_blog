package output

import "time"

// BaseArticle 文章基础信息。
type BaseArticle struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Excerpt       string     `json:"excerpt"`
	FeaturedImage string     `json:"featured_image"`
	AuthorUUID    string     `json:"author_uuid"`
	Status        string     `json:"status"`
	ReadTime      string     `json:"read_time"`
	Views         int32      `json:"views"`
	IsFeatured    bool       `json:"is_featured"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// AuthorInfo 作者信息。
type AuthorInfo struct {
	UUID     string `json:"uuid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// LikeInfo 点赞信息。
type LikeInfo struct {
	Liked *bool `json:"liked,omitempty"`
	Likes int32 `json:"likes"`
}

// BaseCategory 分类基础信息。
type BaseCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// BaseTag 标签基础信息。
type BaseTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ArticleSummary 文章摘要。
type ArticleSummary struct {
	BaseArticle
	Author   AuthorInfo   `json:"author"`
	Like     LikeInfo     `json:"like"`
	Category BaseCategory `json:"category"`
	Tags     []BaseTag    `json:"tags"`
}

// ArticleDetail 文章详情。
type ArticleDetail struct {
	BaseArticle
	Author          AuthorInfo   `json:"author"`
	Like            LikeInfo     `json:"like"`
	Category        BaseCategory `json:"category"`
	Tags            []BaseTag    `json:"tags"`
	Content         string       `json:"content"`
	MetaTitle       string       `json:"meta_title"`
	MetaDescription string       `json:"meta_description"`
}

// ==================== 分类 ====================

// CategoryDetail 分类详情。
type CategoryDetail struct {
	BaseCategory
	ArticleCount int32     `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ==================== 标签 ====================

// TagDetail 标签详情。
type TagDetail struct {
	BaseTag
	ArticleCount int32     `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
