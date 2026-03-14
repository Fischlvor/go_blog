package input

import "io"

// UploadFile 上传文件参数。
type UploadFile struct {
	File        io.Reader
	Filename    string
	Size        int64
	ContentType string
	Usage       string // post_cover, post_content, avatar
	UserUUID    string
}

// ListFiles 文件列表参数。
type ListFiles struct {
	PageParams
	Filename string
	MimeType string
}
