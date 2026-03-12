package entity

import "time"

// Category 分类实体。
type Category struct {
	ID        int64
	Name      string
	Slug      string
	ArticleCount int32
	CreatedAt time.Time
	UpdatedAt time.Time
}
