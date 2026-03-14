package output

import "time"

// UploadResult 上传结果。
type UploadResult struct {
	URL string `json:"url"` // 完整 URL（包含 CDN 域名），前端直接使用
}

// FileInfo 文件信息。
type FileInfo struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}
