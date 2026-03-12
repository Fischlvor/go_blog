package urlutil

import (
	"strings"

	"server-blog-v2/config"
)

// ResolveImageURL 将数据库中的相对路径转换为完整 URL。
// 如果已是完整 URL 则直接返回。
func ResolveImageURL(cfg *config.Config, path string) string {
	if path == "" {
		return ""
	}
	// 已是完整 URL 则直接返回
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	// 拼接七牛域名
	domain := cfg.Qiniu.Domain
	if domain == "" {
		return path
	}
	// 添加协议前缀
	var scheme string
	if cfg.Qiniu.UseHTTPS {
		scheme = "https://"
	} else {
		scheme = "http://"
	}
	// 确保域名以 / 结尾
	if !strings.HasSuffix(domain, "/") {
		domain += "/"
	}
	return scheme + domain + path
}

// ResolveImageURLPtr 处理指针类型的图片路径。
func ResolveImageURLPtr(cfg *config.Config, path *string) string {
	if path == nil {
		return ""
	}
	return ResolveImageURL(cfg, *path)
}
