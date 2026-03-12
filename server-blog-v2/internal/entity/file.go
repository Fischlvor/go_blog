package entity

import "time"

const (
	FileUsagePostCover   = "post_cover"
	FileUsagePostContent = "post_content"
	FileUsageAvatar      = "avatar"
)

// File 文件实体。
type File struct {
	ID         int64
	Key        string
	Filename   string
	Size       int64
	MimeType   string
	Usage      string
	ResourceID *int64
	CreatedAt  time.Time
}
