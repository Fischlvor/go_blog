package request

// BatchDelete 批量删除请求。
type BatchDelete struct {
	IDs []int64 `json:"ids" validate:"required,min=1"`
}

// BatchDeleteStr 批量删除请求（字符串 ID）。
type BatchDeleteStr struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}
