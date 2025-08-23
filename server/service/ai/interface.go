package ai

import (
	"context"
	"io"
)

// Message 聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content      string `json:"content"`
	Tokens       int    `json:"tokens"`
	FinishReason string `json:"finish_reason"`
}

// StreamingChatResponse 流式聊天响应
type StreamingChatResponse struct {
	Content      string `json:"content"`
	IsComplete   bool   `json:"is_complete"`
	Tokens       int    `json:"tokens"`
	FinishReason string `json:"finish_reason"`
}

// AIModelProvider AI模型提供商接口
type AIModelProvider interface {
	// Chat 普通聊天
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// ChatStream 流式聊天
	ChatStream(ctx context.Context, req ChatRequest) (io.ReadCloser, error)

	// GetName 获取提供商名称
	GetName() string

	// IsAvailable 检查是否可用
	IsAvailable() bool
}

// ModelConfig 模型配置
type ModelConfig struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Provider    string  `json:"provider"`
	Endpoint    string  `json:"endpoint"`
	ApiKey      string  `json:"api_key"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	IsActive    bool    `json:"is_active"`
}
