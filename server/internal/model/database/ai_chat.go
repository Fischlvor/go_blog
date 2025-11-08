package database

import (
	"server/pkg/global"

	"github.com/gofrs/uuid"
)

// AIChatSession AI聊天会话
type AIChatSession struct {
	global.MODEL
	UserUUID uuid.UUID `json:"user_uuid" gorm:"type:char(36);not null;comment:用户UUID"`
	Title    string    `json:"title" gorm:"size:255;comment:会话标题"`
	Model    string    `json:"model" gorm:"size:50;not null;comment:AI模型名称"`

	// 关联关系
	User     User            `json:"user" gorm:"foreignKey:UserUUID;references:UUID"` // 用户信息
	Messages []AIChatMessage `json:"messages" gorm:"foreignKey:SessionID"`
}

// AIChatMessage AI聊天消息
type AIChatMessage struct {
	global.MODEL
	SessionID uint   `json:"session_id" gorm:"not null;comment:会话ID"`
	Role      string `json:"role" gorm:"size:20;not null;comment:消息角色(user/assistant)"`
	Content   string `json:"content" gorm:"type:text;not null;comment:消息内容"`
	Tokens    int    `json:"tokens" gorm:"default:0;comment:token数量"`
}

// AIModel AI模型配置
type AIModel struct {
	global.MODEL
	Name        string  `json:"name" gorm:"size:50;unique;not null;comment:模型名称"`
	DisplayName string  `json:"display_name" gorm:"size:100;not null;comment:显示名称"`
	Provider    string  `json:"provider" gorm:"size:50;not null;comment:提供商"`
	Endpoint    string  `json:"endpoint" gorm:"size:255;comment:API端点"`
	ApiKey      string  `json:"api_key" gorm:"size:255;comment:API密钥"`
	MaxTokens   int     `json:"max_tokens" gorm:"default:4096;comment:最大token数"`
	Temperature float64 `json:"temperature" gorm:"default:0.7;comment:温度参数"`
	IsActive    bool    `json:"is_active" gorm:"default:true;comment:是否启用"`
}
