package resource

import "io"

// ResourceUploader 资源上传接口
// 支持切换不同云存储提供商
type ResourceUploader interface {
	// === 普通上传（小文件 <10MB）===
	// Upload 直接上传文件到云存储
	// fileHash 用于生成文件Key，便于秒传检测
	Upload(data []byte, fileName string, mimeType string, fileHash string) (fileKey string, err error)

	// === 分片上传（大文件 ≥10MB）===
	// UploadBlockStream 流式上传块到云存储，返回七牛云Context
	// 优点：不需要将整个块读入内存，边接收边转发，大幅降低内存占用
	UploadBlockStream(reader io.Reader, blockSize int64) (context string, err error)

	// MergeBlocks 合并所有块为最终文件
	MergeBlocks(fileSize int64, fileKey string, contexts []string) error

	// === 通用方法 ===
	// GetFileURL 获取文件的CDN URL
	GetFileURL(fileKey string) string

	// DeleteFile 删除文件
	DeleteFile(fileKey string) error

	// GenerateFileKey 生成文件存储Key
	GenerateFileKey(fileName string, fileHash string) string
}

// FileValidator 文件验证器接口
type FileValidator interface {
	// ValidateMimeType 验证MIME类型
	ValidateMimeType(mimeType string) error

	// ValidateSignature 验证文件签名（Magic Number）
	ValidateSignature(data []byte, mimeType string) error

	// ValidateSize 验证文件大小
	ValidateSize(size int64) error
}
