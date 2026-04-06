package entity

import "time"

// SiteSetting 站点配置实体。
type SiteSetting struct {
	ID          int64
	SettingKey  string
	SettingValue string
	SettingType string
	Description *string
	IsPublic    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
