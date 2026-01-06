package response

// ResourceCheckResponse 检查文件响应
type ResourceCheckResponse struct {
	Exists         bool   `json:"exists"`                    // 是否已存在（秒传成功）
	FileURL        string `json:"file_url,omitempty"`        // 秒传时返回URL
	TaskID         string `json:"task_id,omitempty"`         // 续传时返回任务ID
	TotalChunks    int    `json:"total_chunks,omitempty"`    // 总块数
	UploadedChunks []int  `json:"uploaded_chunks,omitempty"` // 已上传的块号
	MissingChunks  []int  `json:"missing_chunks,omitempty"`  // 缺失的块号
}

// ResourceInitResponse 初始化任务响应
type ResourceInitResponse struct {
	TaskID      string `json:"task_id"`      // 任务ID
	TotalChunks int    `json:"total_chunks"` // 总块数
	ChunkSize   int    `json:"chunk_size"`   // 块大小（字节）
}

// ResourceUploadChunkResponse 上传分片响应
type ResourceUploadChunkResponse struct {
	Success     bool `json:"success"`      // 是否成功
	ChunkNumber int  `json:"chunk_number"` // 块号
}

// ResourceCompleteResponse 完成上传响应
type ResourceCompleteResponse struct {
	FileURL string `json:"file_url"` // 文件CDN URL
	FileKey string `json:"file_key"` // 文件Key
}

// ResourceProgressResponse 上传进度响应
type ResourceProgressResponse struct {
	TaskID         string `json:"task_id"`         // 任务ID
	TotalChunks    int    `json:"total_chunks"`    // 总块数
	UploadedChunks []int  `json:"uploaded_chunks"` // 已上传的块号
	MissingChunks  []int  `json:"missing_chunks"`  // 缺失的块号
	Progress       int    `json:"progress"`        // 进度百分比
}

// ResourceItem 资源列表项
type ResourceItem struct {
	ID              uint   `json:"id"`
	FileKey         string `json:"file_key"`
	FileName        string `json:"file_name"`
	FileURL         string `json:"file_url"`
	FileSize        int64  `json:"file_size"`
	MimeType        string `json:"mime_type"`
	TranscodeStatus int8   `json:"transcode_status"` // 转码状态：0=无需转码, 1=转码中, 2=成功, 3=失败
	TranscodeURL    string `json:"transcode_url"`    // 转码后视频URL
	ThumbnailURL    string `json:"thumbnail_url"`    // 缩略图URL
}
