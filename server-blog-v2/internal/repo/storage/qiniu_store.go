package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"server-blog-v2/internal/repo"
)

const (
	// ResourcePathPrefix 资源存储路径前缀
	ResourcePathPrefix = "resource"
)

type qiniuStore struct {
	mac        *qbox.Mac
	bucket     string
	domain     string
	useHTTPS   bool
	pathPrefix string
	zone       string
}

// NewQiniuStore 创建七牛云存储。
func NewQiniuStore(accessKey, secretKey, bucket, domain, zone string, useHTTPS bool, pathPrefix string) repo.ObjectStore {
	mac := qbox.NewMac(accessKey, secretKey)
	return &qiniuStore{
		mac:        mac,
		bucket:     bucket,
		domain:     domain,
		useHTTPS:   useHTTPS,
		pathPrefix: pathPrefix,
		zone:       zone,
	}
}

func (s *qiniuStore) Upload(ctx context.Context, key string, data io.Reader, size int64, contentType string) (string, error) {
	// 生成上传凭证
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", s.bucket, key),
	}
	upToken := putPolicy.UploadToken(s.mac)

	// 配置
	cfg := storage.Config{
		UseHTTPS:      s.useHTTPS,
		UseCdnDomains: true,
	}

	// 表单上传
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		MimeType: contentType,
	}

	err := formUploader.Put(ctx, &ret, upToken, key, data, size, &putExtra)
	if err != nil {
		return "", fmt.Errorf("qiniu upload error: %w", err)
	}

	return s.GetURL(ret.Key), nil
}

func (s *qiniuStore) Delete(ctx context.Context, key string) error {
	cfg := storage.Config{
		UseHTTPS: s.useHTTPS,
	}
	bucketManager := storage.NewBucketManager(s.mac, &cfg)
	return bucketManager.Delete(s.bucket, key)
}

func (s *qiniuStore) GetURL(key string) string {
	scheme := "http"
	if s.useHTTPS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s", scheme, s.domain, key)
}

// UploadBlock 流式上传块到七牛云，返回 Context。
func (s *qiniuStore) UploadBlock(ctx context.Context, data io.Reader, size int64) (string, error) {
	putPolicy := storage.PutPolicy{Scope: s.bucket}
	upToken := putPolicy.UploadToken(s.mac)

	upHost := s.getUpHost()
	url := fmt.Sprintf("%s/mkblk/%d", upHost, size)

	req, err := http.NewRequestWithContext(ctx, "POST", url, data)
	if err != nil {
		return "", fmt.Errorf("create request error: %w", err)
	}

	req.ContentLength = size
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "UpToken "+upToken)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload block error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload block failed, status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response error: %w", err)
	}

	ctxStr := extractJSONField(string(body), "ctx")
	if ctxStr == "" {
		return "", fmt.Errorf("parse ctx failed, body: %s", string(body))
	}

	return ctxStr, nil
}

// MergeBlocks 合并所有块为最终文件。
func (s *qiniuStore) MergeBlocks(ctx context.Context, fileSize int64, fileKey string, contexts []string) error {
	putPolicy := storage.PutPolicy{Scope: s.bucket}
	upToken := putPolicy.UploadToken(s.mac)

	upHost := s.getUpHost()
	encodedKey := base64.URLEncoding.EncodeToString([]byte(fileKey))
	url := fmt.Sprintf("%s/mkfile/%d/key/%s", upHost, fileSize, encodedKey)

	body := strings.Join(contexts, ",")
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "UpToken "+upToken)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("merge blocks error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("merge blocks failed, status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// GenerateFileKey 生成文件存储 Key。
func (s *qiniuStore) GenerateFileKey(fileName, fileHash string) string {
	ext := filepath.Ext(fileName)
	if fileHash != "" {
		return fmt.Sprintf("%s/%s%s", ResourcePathPrefix, fileHash, ext)
	}
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s/%s%s", ResourcePathPrefix, timestamp, ext)
}

// getUpHost 获取上传域名。
func (s *qiniuStore) getUpHost() string {
	zoneCode := s.getZoneCode()
	if s.useHTTPS {
		return "https://up-" + zoneCode + ".qiniup.com"
	}
	return "http://up-" + zoneCode + ".qiniup.com"
}

// getZoneCode 获取区域代码。
func (s *qiniuStore) getZoneCode() string {
	switch s.zone {
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
		return "z2"
	}
}

// extractJSONField 简单提取 JSON 字段值。
func extractJSONField(jsonStr, field string) string {
	key := fmt.Sprintf(`"%s"`, field)
	idx := strings.Index(jsonStr, key)
	if idx == -1 {
		return ""
	}

	start := idx + len(key)
	for start < len(jsonStr) && (jsonStr[start] == ':' || jsonStr[start] == ' ') {
		start++
	}

	if start >= len(jsonStr) || jsonStr[start] != '"' {
		return ""
	}

	start++
	end := start
	for end < len(jsonStr) && jsonStr[end] != '"' {
		if jsonStr[end] == '\\' && end+1 < len(jsonStr) {
			end += 2
		} else {
			end++
		}
	}

	return jsonStr[start:end]
}
