package output

import "time"

// AIModelInfo AI 模型信息。
type AIModelInfo struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Provider    string    `json:"provider"`
	Endpoint    string    `json:"endpoint"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
