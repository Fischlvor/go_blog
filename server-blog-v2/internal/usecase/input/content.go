package input

// ListArticles 文章列表参数（管理端）。
type ListArticles struct {
	PageParams
	Keyword    *KeywordParams
	Sort       *SortParams
	CategoryID IntFilterParam
	TagID      IntFilterParam
	Status     *string
	Visibility *string // 可见性筛选：public, private
	IsFeatured *bool
}

// ListPublicArticles 文章列表参数（公开端）。
type ListPublicArticles struct {
	PageParams
	Keyword    *KeywordParams
	Sort       *SortParams
	CategoryID IntFilterParam
	TagID      IntFilterParam
}

// ListUserLikedArticles 用户点赞文章列表参数。
type ListUserLikedArticles struct {
	PageParams
}

// CreateArticle 创建文章参数。
type CreateArticle struct {
	Title         string
	Slug          string
	Excerpt       *string
	Content       string
	FeaturedImage *string
	AuthorUUID    string
	CategoryID    int64
	TagIDs        []int64
	Status        string
	Visibility    string // 可见性：public, private
	IsFeatured    bool
}

// UpdateArticle 更新文章参数。
type UpdateArticle struct {
	Slug          string // 用 slug 作为文章标识
	Title         string
	Excerpt       *string
	Content       string
	FeaturedImage *string
	AuthorUUID    string
	CategoryID    int64
	TagIDs        []int64
	Status        string
	Visibility    string // 可见性：public, private
	IsFeatured    bool
}

// ==================== 分类 ====================

// ListCategories 分类列表参数。
type ListCategories struct {
	PageParams
	Keyword *KeywordParams
	Sort    *SortParams
}

// CreateCategory 创建分类参数。
type CreateCategory struct {
	Name string
	Slug string
}

// UpdateCategory 更新分类参数。
type UpdateCategory struct {
	ID   int64
	Name string
	Slug string
}

// ==================== 标签 ====================

// ListTags 标签列表参数。
type ListTags struct {
	PageParams
	Keyword *KeywordParams
	Sort    *SortParams
}

// CreateTag 创建标签参数。
type CreateTag struct {
	Name string
	Slug string
}

// UpdateTag 更新标签参数。
type UpdateTag struct {
	ID   int64
	Name string
	Slug string
}
