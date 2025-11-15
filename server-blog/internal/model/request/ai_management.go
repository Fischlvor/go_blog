package request

// AIModelListRequest AI模型列表查询请求
type AIModelListRequest struct {
	PageInfo
	Name     string `json:"name" form:"name"`         // 模型名称
	Provider string `json:"provider" form:"provider"` // 提供商
}

// AISessionListRequest AI会话列表查询请求
type AISessionListRequest struct {
	PageInfo
	UserUUID string `json:"user_uuid" form:"user_uuid"` // 用户UUID
	Model    string `json:"model" form:"model"`         // 模型名称
}

// AIMessageListRequest AI消息列表查询请求
type AIMessageListRequest struct {
	PageInfo
	SessionID string `json:"session_id" form:"session_id"` // 会话ID
	Role      string `json:"role" form:"role"`             // 消息角色
}
