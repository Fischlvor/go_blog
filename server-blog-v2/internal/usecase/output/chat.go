package output

import "time"

// Session 聊天会话。
type Session struct {
	ID        int64
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Message 聊天消息。
type Message struct {
	ID        int64
	Role      string // user, assistant
	Content   string
	CreatedAt time.Time
}

// SessionAdmin 聊天会话（管理端）。
type SessionAdmin struct {
	ID           int64         `json:"id"`
	Title        string        `json:"title"`
	UserUUID     string        `json:"user_uuid"`
	User         *SessionUser  `json:"user,omitempty"`
	Model        string        `json:"model"`
	MessageCount int           `json:"message_count"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// SessionUser 会话用户信息。
type SessionUser struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// MessageAdmin 消息（管理端）。
type MessageAdmin struct {
	ID        int64     `json:"id"`
	SessionID int64     `json:"session_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Tokens    int       `json:"tokens"`
	CreatedAt time.Time `json:"created_at"`
}
