package request

// CreateTag 创建标签请求。
type CreateTag struct {
	Name string `json:"name" validate:"required,max=50"`
	Slug string `json:"slug" validate:"required,max=50"`
}

// UpdateTag 更新标签请求。
type UpdateTag struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,max=50"`
	Slug string `json:"slug" validate:"required,max=50"`
}
