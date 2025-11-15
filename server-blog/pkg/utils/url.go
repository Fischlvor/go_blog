package utils

import (
	"server/pkg/global"
	"strings"
)

// PublicURLFromDB 将数据库中的 URL（qiniu: key，本地: 相对路径）转换为可对外访问的完整 URL
func PublicURLFromDB(u string) string {
	if u == "" {
		return ""
	}
	// 已是完整 URL 则直接返回
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	// qiniu: 需要拼接域名
	if global.Config.System.OssType == "qiniu" {
		return global.Config.Qiniu.ImgPath + u
	}
	// local: 直接返回相对路径
	return u
}

// PublicURLSlice 批量转换 URL
func PublicURLSlice(urls []string) []string {
	if len(urls) == 0 {
		return urls
	}
	out := make([]string, len(urls))
	for i := range urls {
		out[i] = PublicURLFromDB(urls[i])
	}
	return out
}

// DBURLFromPublic 将外部全量 URL 转换为数据库存储格式（qiniu: 去掉域名保留 key；local: 返回相对路径）
func DBURLFromPublic(u string) string {
	if u == "" {
		return ""
	}
	// 若已是相对路径或 key，直接返回
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return u
	}
	// qiniu: 去掉 ImgPath 前缀
	if global.Config.System.OssType == "qiniu" {
		prefix := global.Config.Qiniu.ImgPath
		if strings.HasPrefix(u, prefix) {
			return strings.TrimPrefix(u, prefix)
		}
	}
	// 其他情况返回原值（避免误删）；如果是 local 且前端意外传了全量，也按原样
	return u
}
