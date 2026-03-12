package request

// CreateCategory 创建分类请求。
type CreateCategory struct {
	Name string `json:"name" validate:"required,max=50"`
	Slug string `json:"slug" validate:"required,max=50"`
}

// UpdateCategory 更新分类请求。
type UpdateCategory struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,max=50"`
	Slug string `json:"slug" validate:"required,max=50"`
}
