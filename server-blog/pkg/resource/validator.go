package resource

import (
	"errors"
	"strings"
)

// DefaultValidator 默认文件验证器
type DefaultValidator struct {
	maxSize      int64             // 最大文件大小（字节）
	allowedMimes map[string]bool   // 允许的MIME类型
	magicNumbers map[string][]byte // 文件签名（Magic Number）
}

// NewDefaultValidator 创建默认验证器
func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{
		maxSize: 500 * 1024 * 1024, // 默认500MB
		allowedMimes: map[string]bool{
			// 图片
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
			// 视频
			"video/mp4":  true,
			"video/webm": true,
			// 文档
			"application/pdf":    true,
			"application/msword": true,
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
			"application/vnd.ms-excel": true,
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
			"application/vnd.ms-powerpoint":                                             true,
			"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
			// 压缩包
			"application/zip":              true,
			"application/x-rar-compressed": true,
			"application/x-7z-compressed":  true,
		},
		magicNumbers: map[string][]byte{
			"image/jpeg":      {0xFF, 0xD8, 0xFF},
			"image/png":       {0x89, 0x50, 0x4E, 0x47},
			"image/gif":       {0x47, 0x49, 0x46, 0x38},
			"image/webp":      {0x52, 0x49, 0x46, 0x46}, // RIFF
			"video/mp4":       {0x00, 0x00, 0x00},       // ftyp (前3字节)
			"application/pdf": {0x25, 0x50, 0x44, 0x46}, // %PDF
			"application/zip": {0x50, 0x4B, 0x03, 0x04},
		},
	}
}

// ValidateMimeType 验证MIME类型
func (v *DefaultValidator) ValidateMimeType(mimeType string) error {
	// 处理带参数的MIME类型，如 "text/plain; charset=utf-8"
	mimeType = strings.Split(mimeType, ";")[0]
	mimeType = strings.TrimSpace(mimeType)

	if !v.allowedMimes[mimeType] {
		return errors.New("不支持的文件类型: " + mimeType)
	}
	return nil
}

// ValidateSignature 验证文件签名
func (v *DefaultValidator) ValidateSignature(data []byte, mimeType string) error {
	magic, exists := v.magicNumbers[mimeType]
	if !exists {
		// 没有定义签名的类型，跳过验证
		return nil
	}

	if len(data) < len(magic) {
		return errors.New("文件数据太短，无法验证签名")
	}

	for i, b := range magic {
		if data[i] != b {
			return errors.New("文件签名不匹配，可能是伪造的文件类型")
		}
	}

	return nil
}

// ValidateSize 验证文件大小
func (v *DefaultValidator) ValidateSize(size int64) error {
	if size <= 0 {
		return errors.New("文件大小无效")
	}
	if size > v.maxSize {
		return errors.New("文件大小超过限制")
	}
	return nil
}

// SetMaxSize 设置最大文件大小
func (v *DefaultValidator) SetMaxSize(size int64) {
	v.maxSize = size
}

// AddAllowedMime 添加允许的MIME类型
func (v *DefaultValidator) AddAllowedMime(mimeType string) {
	v.allowedMimes[mimeType] = true
}
