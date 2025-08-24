package response

import (
	"time"
)

// 事件类型常量
const (
	EventMessage        = 1 // 正常消息
	EventComplete       = 2 // 流式响应完成
	EventTitleGenerated = 3 // 标题生成完成
)

// ChatSessionResponse 聊天会话响应
type ChatSessionResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Model     string    `json:"model"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatMessageResponse 聊天消息响应
type ChatMessageResponse struct {
	ID        uint      `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Tokens    int       `json:"tokens"`
	CreatedAt time.Time `json:"created_at"`
}

// AIModelResponse AI模型响应
type AIModelResponse struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Provider    string  `json:"provider"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	IsActive    bool    `json:"is_active"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Message     ChatMessageResponse `json:"message"`
	SessionID   uint                `json:"session_id"`
	TotalTokens int                 `json:"total_tokens"`
}

// StreamingChatResponse 流式聊天响应
type StreamingChatResponse struct {
	Content   string `json:"content"`
	SessionID uint   `json:"session_id"`
	MessageID uint   `json:"message_id"` // 消息ID，用于前端定位和更新消息
	EventID   int    `json:"event_id"`   // 事件ID：1(正常消息)、2(完成)、3(标题生成完成)
}
