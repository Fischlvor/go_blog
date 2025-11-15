package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Emoji 表情模型
type Emoji struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Key             string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"key"`
	Filename        string     `gorm:"type:varchar(255);not null" json:"filename"`
	GroupKey        string     `gorm:"type:varchar(50);not null;index" json:"group_key"`
	SpriteGroup     int        `gorm:"not null;index" json:"sprite_group"`
	SpritePositionX int        `gorm:"not null" json:"sprite_position_x"`
	SpritePositionY int        `gorm:"not null" json:"sprite_position_y"`
	FileSize        int        `gorm:"default:0" json:"file_size"`
	CdnUrl          string     `gorm:"type:varchar(500);default:''" json:"cdn_url"`
	UploadTime      time.Time  `gorm:"autoCreateTime" json:"upload_time"`
	Status          int8       `gorm:"default:1;index" json:"status"` // 1=active, 0=deleted
	CreatedBy       uuid.UUID  `gorm:"type:char(36);index" json:"created_by"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Emoji) TableName() string {
	return "emojis"
}

// EmojiGroup 表情组模型
type EmojiGroup struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupName     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"group_name"`
	GroupKey      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"group_key"` // 用于路径的唯一标识
	Description   string    `gorm:"type:varchar(500);default:''" json:"description"`
	SortOrder     int       `gorm:"default:0;index" json:"sort_order"`
	EmojiCount    int       `gorm:"default:0" json:"emoji_count"`
	Status        int8      `gorm:"default:1;index" json:"status"`                       // 1=active, 0=inactive
	SpriteConfUrl string    `gorm:"type:varchar(500);default:''" json:"sprite_conf_url"` // 雪碧图配置文件URL
	CreatedBy     uuid.UUID `gorm:"type:char(36)" json:"created_by"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (EmojiGroup) TableName() string {
	return "emoji_groups"
}

// EmojiSprite 雪碧图模型
type EmojiSprite struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SpriteGroup int       `gorm:"uniqueIndex;not null" json:"sprite_group"`
	Filename    string    `gorm:"type:varchar(255);not null" json:"filename"`
	CdnUrl      string    `gorm:"type:varchar(500);default:''" json:"cdn_url"`
	Width       int       `gorm:"not null" json:"width"`
	Height      int       `gorm:"not null" json:"height"`
	EmojiCount  int       `gorm:"not null" json:"emoji_count"`
	FileSize    int       `gorm:"default:0" json:"file_size"`
	Status      int8      `gorm:"default:1;index" json:"status"` // 1=active, 0=inactive
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (EmojiSprite) TableName() string {
	return "emoji_sprites"
}

// EmojiTask 任务模型
type EmojiTask struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskType    string     `gorm:"type:varchar(50);not null;index" json:"task_type"`
	Status      string     `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	Progress    int        `gorm:"default:0" json:"progress"`
	Message     string     `gorm:"type:varchar(500);default:''" json:"message"`
	Params      string     `gorm:"type:json" json:"params,omitempty"`
	Result      string     `gorm:"type:json" json:"result,omitempty"`
	CreatedBy   uuid.UUID  `gorm:"type:char(36);index" json:"created_by"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (EmojiTask) TableName() string {
	return "emoji_tasks"
}

// 常量定义
const (
	// Emoji状态
	EmojiStatusActive  = 1
	EmojiStatusDeleted = 0

	// EmojiSprite状态
	EmojiSpriteStatusActive   = 1
	EmojiSpriteStatusInactive = 0

	// 任务状态
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"

	// 任务类型
	TaskTypeRegenerateSprites = "regenerate_sprites"
	TaskTypeUploadEmoji       = "upload_emoji"
	TaskTypeBatchUpload       = "batch_upload"
)

// EmojiListResponse 表情列表响应
type EmojiListResponse struct {
	List     []Emoji `json:"list"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

// EmojiGroupListResponse 表情组列表响应
type EmojiGroupListResponse struct {
	List     []EmojiGroup `json:"list"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

// EmojiSpriteListResponse 雪碧图列表响应
type EmojiSpriteListResponse struct {
	List     []EmojiSprite `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// EmojiUploadRequest 上传请求
type EmojiUploadRequest struct {
	GroupKey string `json:"group_key" binding:"required"`
	Files    []struct {
		Filename string `json:"filename"`
		Content  string `json:"content"` // base64编码的文件内容
	} `json:"files" binding:"required"`
}

// EmojiConfigResponse 前端配置响应
type EmojiConfigResponse struct {
	Version     string            `json:"version"`
	TotalEmojis int64             `json:"total_emojis"`
	Sprites     []EmojiSpriteInfo `json:"sprites"`
	Mapping     map[string]string `json:"mapping"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// EmojiSpriteInfo 雪碧图信息
type EmojiSpriteInfo struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Range    [2]int `json:"range"` // [start, end]
	Frozen   bool   `json:"frozen"`
	Size     [2]int `json:"size"` // [width, height]
}
