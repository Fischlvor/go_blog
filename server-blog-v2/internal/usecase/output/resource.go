package output

import "time"

// ResourceCheckResponse 检查文件响应。
type ResourceCheckResponse struct {
	Exists         bool   `json:"exists"`
	FileURL        string `json:"file_url,omitempty"`
	TaskID         string `json:"task_id,omitempty"`
	TotalChunks    int    `json:"total_chunks,omitempty"`
	UploadedChunks []int  `json:"uploaded_chunks,omitempty"`
	MissingChunks  []int  `json:"missing_chunks,omitempty"`
}

// ResourceInitResponse 初始化上传任务响应。
type ResourceInitResponse struct {
	TaskID      string `json:"task_id"`
	TotalChunks int    `json:"total_chunks"`
	ChunkSize   int    `json:"chunk_size"`
}

// ResourceUploadChunkResponse 上传分片响应。
type ResourceUploadChunkResponse struct {
	Success     bool `json:"success"`
	ChunkNumber int  `json:"chunk_number"`
}

// ResourceCompleteResponse 完成上传响应。
type ResourceCompleteResponse struct {
	FileURL string `json:"file_url"`
	FileKey string `json:"file_key"`
}

// ResourceProgressResponse 查询上传进度响应。
type ResourceProgressResponse struct {
	TaskID         string `json:"task_id"`
	TotalChunks    int    `json:"total_chunks"`
	UploadedChunks []int  `json:"uploaded_chunks"`
	MissingChunks  []int  `json:"missing_chunks"`
	Progress       int    `json:"progress"`
}

// ResourceInfo 资源信息。
type ResourceInfo struct {
	ID              int64     `json:"id"`
	FileKey         string    `json:"file_key"`
	FileName        string    `json:"file_name"`
	FileURL         string    `json:"file_url"`
	FileSize        int64     `json:"file_size"`
	MimeType        string    `json:"mime_type"`
	TranscodeStatus int8      `json:"transcode_status"`
	TranscodeURL    *string   `json:"transcode_url,omitempty"`
	ThumbnailURL    *string   `json:"thumbnail_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
