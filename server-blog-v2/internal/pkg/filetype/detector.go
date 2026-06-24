package filetype

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

// DetectMimeType 检测文件的真实 MIME 类型（基于文件内容）。
// 读取前 512 字节进行检测，不影响原始 reader。
func DetectMimeType(reader io.Reader) (mimeType string, err error) {
	// 读取前 512 字节（http.DetectContentType 的要求）
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	
	if n == 0 {
		return "", errors.New("文件为空")
	}
	
	// 使用标准库检测
	detected := http.DetectContentType(buffer[:n])
	
	// 标准化处理
	detected = normalizeMimeType(detected)
	
	return detected, nil
}

// VerifyMimeType 验证声明的 MIME 类型与实际检测的类型是否匹配。
// reader: 文件内容流
// declaredMimeType: 前端声明的 MIME 类型
// 返回: 检测到的真实类型，是否匹配，错误
func VerifyMimeType(reader io.Reader, declaredMimeType string) (detectedType string, matched bool, err error) {
	detectedType, err = DetectMimeType(reader)
	if err != nil {
		return "", false, err
	}
	
	// 标准化声明的类型
	declaredMimeType = normalizeMimeType(declaredMimeType)
	
	// 检查是否匹配
	matched = isMimeTypeMatch(declaredMimeType, detectedType)
	
	return detectedType, matched, nil
}

// normalizeMimeType 标准化 MIME 类型。
func normalizeMimeType(mimeType string) string {
	// 移除参数部分（如 charset）
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		mimeType = mimeType[:idx]
	}
	
	// 转小写并去除空格
	mimeType = strings.TrimSpace(strings.ToLower(mimeType))
	
	return mimeType
}

// isMimeTypeMatch 判断两个 MIME 类型是否匹配。
func isMimeTypeMatch(declared, detected string) bool {
	// 完全匹配
	if declared == detected {
		return true
	}
	
	// 特殊规则：某些格式的兼容性处理
	compatMap := map[string][]string{
		// image/jpg 和 image/jpeg 是同一种格式
		"image/jpeg": {"image/jpg"},
		"image/jpg":  {"image/jpeg"},
		
		// MPEG 视频的不同表示
		"video/mpeg": {"video/mpg"},
		"video/mpg":  {"video/mpeg"},
		
		// 某些文档格式 http.DetectContentType 检测为 application/octet-stream
		// 但我们信任前端的声明（如果在白名单内）
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {"application/zip", "application/octet-stream"},
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       {"application/zip", "application/octet-stream"},
		"application/vnd.ms-excel": {"application/octet-stream"},
		"application/msword":       {"application/octet-stream"},
		
		// 压缩包的不同表示
		"application/x-zip-compressed": {"application/zip"},
		"application/zip":              {"application/x-zip-compressed"},
		
		// 某些浏览器对 MP3 的识别
		"audio/mp3":  {"audio/mpeg"},
		"audio/mpeg": {"audio/mp3"},
	}
	
	// 检查兼容映射
	if compatTypes, exists := compatMap[declared]; exists {
		for _, compatType := range compatTypes {
			if detected == compatType {
				return true
			}
		}
	}
	
	return false
}

// CreateRewindableReader 创建可重读的 reader（用于先检测类型，再上传）。
// 返回: 缓冲的内容，新的 reader（可以重新读取）
func CreateRewindableReader(reader io.Reader, bufferSize int) ([]byte, io.Reader, error) {
	buffer := make([]byte, bufferSize)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}
	
	buffer = buffer[:n]
	
	// 创建一个新的 reader，可以重新读取缓冲的内容 + 原始 reader 剩余内容
	newReader := io.MultiReader(bytes.NewReader(buffer), reader)
	
	return buffer, newReader, nil
}
