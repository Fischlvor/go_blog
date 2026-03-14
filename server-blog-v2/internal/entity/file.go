package entity

import "time"

const (
	FileUsageArticleCover   = "article_cover"
	FileUsageArticleContent = "article_content"
	FileUsageAvatar         = "avatar"
)

// File 文件实体。
type File struct {
	ID         int64
	Key        string
	Filename   string
	FileHash   string
	Size       int64
	MimeType   string
	Usage      string
	ResourceID *int64
	UserUUID   string
	CreatedAt  time.Time
}
