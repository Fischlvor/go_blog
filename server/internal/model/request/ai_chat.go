package request

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	Title string `json:"title" binding:"required" msg:"会话标题不能为空"`
	Model string `json:"model" binding:"required" msg:"AI模型不能为空"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	SessionID uint   `json:"session_id" binding:"required" msg:"会话ID不能为空"`
	Content   string `json:"content" binding:"required" msg:"消息内容不能为空"`
}

// GetSessionsRequest 获取会话列表请求
type GetSessionsRequest struct {
	PageInfo
}

// GetMessagesRequest 获取消息列表请求
type GetMessagesRequest struct {
	SessionID uint `json:"session_id" binding:"required" msg:"会话ID不能为空"`
	PageInfo
}

// DeleteSessionRequest 删除会话请求
type DeleteSessionRequest struct {
	SessionID uint `json:"session_id" binding:"required" msg:"会话ID不能为空"`
}

// UpdateSessionRequest 更新会话请求
type UpdateSessionRequest struct {
	SessionID uint   `json:"session_id" binding:"required" msg:"会话ID不能为空"`
	Title     string `json:"title" binding:"required" msg:"会话标题不能为空"`
}
