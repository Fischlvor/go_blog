package output

import "time"

// EmojiConfig 表情配置。
type EmojiConfig struct {
	Version     string            `json:"version"`
	TotalEmojis int64             `json:"total_emojis"`
	Sprites     []EmojiSpriteInfo `json:"sprites"`
	Mapping     map[string]string `json:"mapping"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// EmojiSpriteInfo 雪碧图信息。
type EmojiSpriteInfo struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Range    [2]int `json:"range"`
	Frozen   bool   `json:"frozen"`
	Size     [2]int `json:"size"`
}

// EmojiGroup 表情组。
type EmojiGroup struct {
	ID          int64     `json:"id"`
	GroupName   string    `json:"group_name"`
	GroupKey    string    `json:"group_key"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	EmojiCount  int       `json:"emoji_count"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EmojiInfo 表情信息。
type EmojiInfo struct {
	ID              int64     `json:"id"`
	Key             string    `json:"key"`
	Filename        string    `json:"filename"`
	GroupKey        string    `json:"group_key"`
	SpriteGroup     int       `json:"sprite_group"`
	SpritePositionX int       `json:"sprite_position_x"`
	SpritePositionY int       `json:"sprite_position_y"`
	FileSize        int64     `json:"file_size"`
	CdnURL          string    `json:"cdn_url"`
	UploadTime      time.Time `json:"upload_time"`
	Status          int       `json:"status"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SpriteInfo 雪碧图信息（管理端）。
type SpriteInfo struct {
	ID          int64     `json:"id"`
	SpriteGroup int       `json:"sprite_group"`
	Filename    string    `json:"filename"`
	CdnURL      string    `json:"cdn_url"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	EmojiCount  int       `json:"emoji_count"`
	FileSize    int       `json:"file_size"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
