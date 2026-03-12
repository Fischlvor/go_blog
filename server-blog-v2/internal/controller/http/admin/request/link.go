package request

// CreateLink 创建友链请求。
type CreateLink struct {
	Name        string `json:"name" validate:"required,max=100"`
	URL         string `json:"url" validate:"required,url,max=500"`
	Logo        string `json:"logo" validate:"max=500"`
	Description string `json:"description" validate:"max=200"`
	Sort        int    `json:"sort"`
}

// UpdateLink 更新友链请求。
type UpdateLink struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,max=100"`
	URL         string `json:"url" validate:"required,url,max=500"`
	Logo        string `json:"logo" validate:"max=500"`
	Description string `json:"description" validate:"max=200"`
	Sort        int    `json:"sort"`
}
