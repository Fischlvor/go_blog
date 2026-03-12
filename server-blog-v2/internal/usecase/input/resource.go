package input

import "io"

// ResourceCheck 检查文件请求。
type ResourceCheck struct {
	FileHash string `json:"file_hash"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

// ResourceInit 初始化上传任务请求。
type ResourceInit struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileHash string `json:"file_hash"`
	MimeType string `json:"mime_type"`
}

// ResourceUploadChunk 上传分片请求。
type ResourceUploadChunk struct {
	TaskID      string
	ChunkNumber int
	ChunkData   io.Reader
	ChunkSize   int64
}

// ResourceComplete 完成上传请求。
type ResourceComplete struct {
	TaskID string `json:"task_id"`
}

// ResourceCancel 取消上传请求。
type ResourceCancel struct {
	TaskID string `json:"task_id"`
}

// ResourceProgress 查询上传进度请求。
type ResourceProgress struct {
	TaskID string `query:"task_id"`
}

// ListResources 资源列表请求。
type ListResources struct {
	PageParams
	Filename string
	MimeType string
}

// QiniuCallbackItem 七牛云回调中的单个处理结果。
type QiniuCallbackItem struct {
	Cmd       string `json:"cmd"`
	Code      int    `json:"code"`
	Desc      string `json:"desc"`
	Key       string `json:"key"`
	Hash      string `json:"hash"`
	Fsize     int64  `json:"fsize"`
	ReturnOld int    `json:"returnOld"`
}

// QiniuCallback 七牛云转码回调请求。
type QiniuCallback struct {
	ID           string              `json:"id"`
	Pipeline     string              `json:"pipeline"`
	Code         int                 `json:"code"`
	Desc         string              `json:"desc"`
	Reqid        string              `json:"reqid"`
	InputBucket  string              `json:"inputBucket"`
	InputKey     string              `json:"inputKey"`
	Items        []QiniuCallbackItem `json:"items"`
	CreationDate string              `json:"creationDate"`
}
