package entity

import "time"

// Emoji 表情实体。
type Emoji struct {
	ID              int64
	Key             string
	Filename        string
	GroupKey        string
	SpriteGroup     int
	SpritePositionX int
	SpritePositionY int
	FileSize        int
	CdnURL          string
	UploadTime      time.Time
	Status          int8
	CreatedBy       string
	DeletedAt       *time.Time
	UpdatedAt       time.Time
}

// EmojiGroup 表情组实体。
type EmojiGroup struct {
	ID            int64
	GroupName     string
	GroupKey      string
	Description   string
	SortOrder     int
	EmojiCount    int
	Status        int8
	SpriteConfURL string
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// EmojiSprite 雪碧图实体。
type EmojiSprite struct {
	ID          int64
	SpriteGroup int
	Filename    string
	CdnURL      string
	Width       int
	Height      int
	EmojiCount  int
	FileSize    int
	Status      int8
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Emoji 状态常量。
const (
	EmojiStatusActive  = 1
	EmojiStatusDeleted = 0

	EmojiSpriteStatusActive   = 1
	EmojiSpriteStatusInactive = 0
)
