package request

// CreateArticle 创建文章请求。
type CreateArticle struct {
	Title         string  `json:"title" validate:"required,max=200"`
	Slug          string  `json:"slug" validate:"omitempty,max=200"` // 可选，为空时后端自动生成
	Content       string  `json:"content" validate:"required"`
	Excerpt       string  `json:"excerpt" validate:"max=500"`
	FeaturedImage string  `json:"featured_image"`
	CategoryID    int64   `json:"category_id" validate:"required"`
	TagIDs        []int64 `json:"tag_ids"`
	Status        string  `json:"status" validate:"required,oneof=draft published"`
	Visibility    string  `json:"visibility" validate:"omitempty,oneof=public private"` // 可见性：public, private
	IsFeatured    bool    `json:"is_featured"`
}

// UpdateArticle 更新文章请求。
type UpdateArticle struct {
	Slug          string  `json:"slug" validate:"required,max=200"` // 用 slug 作为文章标识
	Title         string  `json:"title" validate:"required,max=200"`
	Content       string  `json:"content"` // 可选，为空时保留原内容
	Excerpt       string  `json:"excerpt" validate:"max=500"`
	FeaturedImage string  `json:"featured_image"`
	CategoryID    int64   `json:"category_id" validate:"required"`
	TagIDs        []int64 `json:"tag_ids"`
	Status        string  `json:"status" validate:"required,oneof=draft published"`
	Visibility    string  `json:"visibility" validate:"omitempty,oneof=public private"` // 可见性：public, private
	IsFeatured    bool    `json:"is_featured"`
}
