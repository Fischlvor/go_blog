package input

// CreateLink 创建友链参数。
type CreateLink struct {
	Name        string
	URL         string
	Logo        string
	Description string
	Sort        int
	IsVisible   bool
}

// UpdateLink 更新友链参数。
type UpdateLink struct {
	ID          int64
	Name        string
	URL         string
	Logo        string
	Description string
	Sort        int
	IsVisible   bool
}
