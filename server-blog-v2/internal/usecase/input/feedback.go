package input

// CreateFeedback 创建反馈参数。
type CreateFeedback struct {
	UserUUID string // 可选，登录用户
	Type     string
	Content  string
	Contact  string
}

// ListFeedback 反馈列表参数。
type ListFeedback struct {
	PageParams
	Status *string
}
