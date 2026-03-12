package entity

import "time"

// ArticleView 文章浏览记录。
type ArticleView struct {
	ID          int64
	ArticleSlug string
	IPAddress   string
	UserAgent   *string
	Referer     *string
	ViewedAt    time.Time
}
