package response

import (
	"time"
)

// ChatSessionResponse 聊天会话响应
type ChatSessionResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 最后一条消息预览
	LastMessage  string `json:"last_message"`
	MessageCount int    `json:"message_count"`
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
	Content     string `json:"content"`
	IsComplete  bool   `json:"is_complete"`
	SessionID   uint   `json:"session_id"`
	MessageID   uint   `json:"message_id"`
	TotalTokens int    `json:"total_tokens"`
}
