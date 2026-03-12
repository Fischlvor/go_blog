package request

// CreateArticle 创建文章请求。
type CreateArticle struct {
	Title         string  `json:"title" validate:"required,max=200"`
	Slug          string  `json:"slug" validate:"required,max=200"`
	Content       string  `json:"content" validate:"required"`
	Excerpt       string  `json:"excerpt" validate:"max=500"`
	FeaturedImage string  `json:"featured_image"`
	CategoryID    int64   `json:"category_id" validate:"required"`
	TagIDs        []int64 `json:"tag_ids"`
	Status        string  `json:"status" validate:"required,oneof=draft published"`
	IsFeatured    bool    `json:"is_featured"`
}

// UpdateArticle 更新文章请求。
type UpdateArticle struct {
	ID            int64   `json:"id" validate:"required"`
	Title         string  `json:"title" validate:"required,max=200"`
	Slug          string  `json:"slug" validate:"required,max=200"`
	Content       string  `json:"content" validate:"required"`
	Excerpt       string  `json:"excerpt" validate:"max=500"`
	FeaturedImage string  `json:"featured_image"`
	CategoryID    int64   `json:"category_id" validate:"required"`
	TagIDs        []int64 `json:"tag_ids"`
	Status        string  `json:"status" validate:"required,oneof=draft published"`
	IsFeatured    bool    `json:"is_featured"`
}
