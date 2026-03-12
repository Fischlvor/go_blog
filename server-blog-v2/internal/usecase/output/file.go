package output

import "time"

// UploadResult 上传结果。
type UploadResult struct {
	Key string
	URL string
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

