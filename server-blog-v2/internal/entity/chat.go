package entity

import "time"

// ChatSession AI 聊天会话。
type ChatSession struct {
	ID        int64
	UserUUID  string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChatMessage AI 聊天消息。
type ChatMessage struct {
	ID        int64
	SessionID int64
	Role      string // user, assistant
	Content   string
	Tokens    int
	CreatedAt time.Time
}

// AIModel AI 模型配置。
type AIModel struct {
	ID          int64
	Name        string
	DisplayName string
	Provider    string
	Endpoint    string
	ApiKey      string
	MaxTokens   int
	Temperature float64
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
