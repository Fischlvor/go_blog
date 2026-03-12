package entity

import "time"

// Tag 标签实体。
type Tag struct {
	ID        int64
	Name      string
	Slug      string
	ArticleCount int32
	CreatedAt time.Time
	UpdatedAt time.Time
}
