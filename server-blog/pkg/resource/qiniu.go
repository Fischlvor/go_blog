package resource

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"server/pkg/global"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const (
	// ResourcePathPrefix 资源存储路径前缀
	ResourcePathPrefix = "resource"
)

// QiniuResourceUploader 七牛云资源上传器
type QiniuResourceUploader struct {
	mac    *qbox.Mac
	cfg    *storage.Config
	bucket string
	domain string
}

// NewQiniuResourceUploader 创建七牛云资源上传器
func NewQiniuResourceUploader() *QiniuResourceUploader {
	mac := qbox.NewMac(global.Config.Qiniu.AccessKey, global.Config.Qiniu.SecretKey)
	cfg := qiniuConfig()
	return &QiniuResourceUploader{
		mac:    mac,
		cfg:    cfg,
		bucket: global.Config.Qiniu.Bucket,
		domain: global.Config.Qiniu.ImgPath,
	}
}

// Upload 普通上传（小文件）
// fileHash 用于生成文件Key，便于秒传检测
func (q *QiniuResourceUploader) Upload(data []byte, fileName string, mimeType string, fileHash string) (string, error) {
	putPolicy := storage.PutPolicy{Scope: q.bucket}
	upToken := putPolicy.UploadToken(q.mac)

	formUploader := storage.NewFormUploader(q.cfg)
	putRet := storage.PutRet{}
	putExtra := storage.PutExtra{
		MimeType: mimeType,
	}

	// 生成文件Key
	fileKey := q.GenerateFileKey(fileName, fileHash)

	err := formUploader.Put(
		context.Background(),
		&putRet,
		upToken,
		fileKey,
		bytes.NewReader(data),
		int64(len(data)),
		&putExtra,
	)
	if err != nil {
		return "", fmt.Errorf("七牛云上传失败: %w", err)
	}

	return putRet.Key, nil
}

// UploadBlockStream 流式上传块到七牛云，返回Context
// 使用七牛云分片上传V1 API: POST /mkblk/<blockSize>
// 优点：不需要将整个块读入内存，边接收边转发
func (q *QiniuResourceUploader) UploadBlockStream(reader io.Reader, blockSize int64) (string, error) {
	putPolicy := storage.PutPolicy{Scope: q.bucket}
	upToken := putPolicy.UploadToken(q.mac)

	// 获取上传域名
	upHost, err := q.getUpHost()
	if err != nil {
		return "", fmt.Errorf("获取上传域名失败: %w", err)
	}

	// 构建 mkblk 请求
	url := fmt.Sprintf("%s/mkblk/%d", upHost, blockSize)

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.ContentLength = blockSize // 必须设置 Content-Length
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "UpToken "+upToken)

	client := &http.Client{Timeout: 120 * time.Second} // 流式上传可能需要更长时间
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传块失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传块失败, 状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应获取 ctx
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 响应格式: {"ctx":"xxx","checksum":"xxx","crc32":xxx,"offset":xxx,"host":"xxx","expired_at":xxx}
	// 简单解析 ctx 字段
	ctxStr := extractJSONField(string(body), "ctx")
	if ctxStr == "" {
		return "", fmt.Errorf("解析ctx失败, 响应: %s", string(body))
	}

	return ctxStr, nil
}

// MergeBlocks 合并所有块为最终文件
// 使用七牛云分片上传V1 API: POST /mkfile/<fileSize>/key/<encodedKey>
func (q *QiniuResourceUploader) MergeBlocks(fileSize int64, fileKey string, contexts []string) error {
	putPolicy := storage.PutPolicy{Scope: q.bucket}
	upToken := putPolicy.UploadToken(q.mac)

	// 获取上传域名
	upHost, err := q.getUpHost()
	if err != nil {
		return fmt.Errorf("获取上传域名失败: %w", err)
	}

	// Base64编码文件Key
	encodedKey := base64.URLEncoding.EncodeToString([]byte(fileKey))

	// 构建 mkfile 请求
	url := fmt.Sprintf("%s/mkfile/%d/key/%s", upHost, fileSize, encodedKey)

	// 请求体：所有context用逗号连接
	body := strings.Join(contexts, ",")

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "UpToken "+upToken)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("合并文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("合并文件失败, 状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// GetFileURL 获取文件的CDN URL
func (q *QiniuResourceUploader) GetFileURL(fileKey string) string {
	domain := strings.TrimSuffix(q.domain, "/")
	return domain + "/" + fileKey
}

// DeleteFile 删除文件
func (q *QiniuResourceUploader) DeleteFile(fileKey string) error {
	bucketManager := storage.NewBucketManager(q.mac, q.cfg)
	return bucketManager.Delete(q.bucket, fileKey)
}

// GenerateFileKey 生成文件存储Key
// fileHash 必须传入，用于秒传检测和文件去重
func (q *QiniuResourceUploader) GenerateFileKey(fileName string, fileHash string) string {
	ext := filepath.Ext(fileName)

	if fileHash != "" {
		// 使用文件hash作为文件名（便于秒传和去重）
		return fmt.Sprintf("%s/%s%s", ResourcePathPrefix, fileHash, ext)
	}

	// 兜底：使用时间戳作为文件名（不推荐，应始终传入hash）
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s/%s%s", ResourcePathPrefix, timestamp, ext)
}

// getUpHost 获取上传域名
func (q *QiniuResourceUploader) getUpHost() (string, error) {
	if q.cfg.UseHTTPS {
		return "https://up-" + q.getZoneCode() + ".qiniup.com", nil
	}
	return "http://up-" + q.getZoneCode() + ".qiniup.com", nil
}

// getZoneCode 获取区域代码
func (q *QiniuResourceUploader) getZoneCode() string {
	switch global.Config.Qiniu.Zone {
	case "z0", "ZoneHuadong":
		return "z0"
	case "z1", "ZoneHuabei":
		return "z1"
	case "z2", "ZoneHuanan":
		return "z2"
	case "na0", "ZoneBeimei":
		return "na0"
	case "as0", "ZoneXinjiapo":
		return "as0"
	default:
		return "z2" // 默认华南
	}
}

// qiniuConfig 创建七牛云配置
func qiniuConfig() *storage.Config {
	cfg := storage.Config{
		UseHTTPS:      global.Config.Qiniu.UseHTTPS,
		UseCdnDomains: global.Config.Qiniu.UseCdnDomains,
	}
	switch global.Config.Qiniu.Zone {
	case "z0", "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "z1", "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "z2", "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "na0", "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "as0", "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	case "ZoneHuadongZheJiang2":
		cfg.Zone = &storage.ZoneHuadongZheJiang2
	}
	return &cfg
}

// extractJSONField 简单提取JSON字段值（避免引入额外依赖）
func extractJSONField(jsonStr, field string) string {
	// 查找 "field":"value" 或 "field": "value"
	key := fmt.Sprintf(`"%s"`, field)
	idx := strings.Index(jsonStr, key)
	if idx == -1 {
		return ""
	}

	// 跳过 key 和冒号
	start := idx + len(key)
	for start < len(jsonStr) && (jsonStr[start] == ':' || jsonStr[start] == ' ') {
		start++
	}

	if start >= len(jsonStr) || jsonStr[start] != '"' {
		return ""
	}

	// 找到值的结束位置
	start++ // 跳过开头的引号
	end := start
	for end < len(jsonStr) && jsonStr[end] != '"' {
		if jsonStr[end] == '\\' && end+1 < len(jsonStr) {
			end += 2 // 跳过转义字符
		} else {
			end++
		}
	}

	return jsonStr[start:end]
}
