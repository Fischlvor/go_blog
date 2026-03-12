package response

import "time"

// ==================== 错误码 ====================

const (
	// 文章
	ErrorListPostsFailed = "0101"
	ErrorGetPostFailed   = "0102"
	ErrorPostNotFound    = "0103"
	ErrorLikePostFailed  = "0104"
	ErrorUnlikePostFailed = "0105"

	// 分类
	ErrorListCategoriesFailed = "0111"

	// 标签
	ErrorListTagsFailed = "0121"

	// 通用
	ErrorParamMissing    = "0001"
	ErrorParamFormat     = "0002"
	ErrorParam           = "0003"
	ErrorLoginRequired   = "0020"
	ErrorTokenInvalid    = "0021"
	ErrorTokenExpired    = "0022"
	ErrorNotImplemented  = "0099"

	// 用户
	ErrorUserNotFound    = "0201"
)

// ==================== 文章响应 ====================

// ArticleSummary 文章摘要响应。
type ArticleSummary struct {
	ID            int64        `json:"id"`
	Title         string       `json:"title"`
	Slug          string       `json:"slug"`
	Excerpt       string       `json:"excerpt"`
	FeaturedImage string       `json:"featured_image"`
	AuthorUUID    string       `json:"author_uuid"`
	Author        AuthorInfo   `json:"author"`
	Status        string       `json:"status"`
	ReadTime      int          `json:"read_time"`
	Views         int32        `json:"views"`
	Like          LikeInfo     `json:"like"`
	IsFeatured    bool         `json:"is_featured"`
	PublishedAt   *time.Time   `json:"published_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Category      BaseCategory `json:"category"`
	Tags          []BaseTag    `json:"tags"`
}

// ArticleDetail 文章详情响应。
type ArticleDetail struct {
	ID              int64        `json:"id"`
	Title           string       `json:"title"`
	Slug            string       `json:"slug"`
	Excerpt         string       `json:"excerpt"`
	FeaturedImage   string       `json:"featured_image"`
	AuthorUUID      string       `json:"author_uuid"`
	Author          AuthorInfo   `json:"author"`
	Status          string       `json:"status"`
	ReadTime        int          `json:"read_time"`
	Views           int32        `json:"views"`
	Like            LikeInfo     `json:"like"`
	IsFeatured      bool         `json:"is_featured"`
	PublishedAt     *time.Time   `json:"published_at"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Category        BaseCategory `json:"category"`
	Tags            []BaseTag    `json:"tags"`
	Content         string       `json:"content"`
	MetaTitle       string       `json:"meta_title"`
	MetaDescription string       `json:"meta_description"`
}

// AuthorInfo 作者信息。
type AuthorInfo struct {
	UUID           string `json:"uuid"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`
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

// CategoryDetail 分类详情响应。
type CategoryDetail struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ArticleCount int32     `json:"article_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TagDetail 标签详情响应。
type TagDetail struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ArticleCount int32     `json:"article_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ArticleSummaryPage 文章摘要分页响应。
type ArticleSummaryPage struct {
	List        []ArticleSummary `json:"list"`
	CurrentPage int           `json:"current_page"`
	PageSize    int           `json:"page_size"`
	TotalItems  int64         `json:"total_items"`
	TotalPages  int           `json:"total_pages"`
}
