package service

import (
	"auth-service/internal/model/database"
	"auth-service/internal/model/response"
	"auth-service/pkg/core"
	"auth-service/pkg/global"
	"strings"

	"go.uber.org/zap"
)

type ApplicationService struct{}

// GetPublicApplications 获取所有公开应用列表
func (s *ApplicationService) GetPublicApplications() ([]response.PublicApplicationInfo, error) {
	var apps []database.SSOApplication

	// 查询所有启用且公开的应用
	err := global.DB.Where("status = ? AND is_public = ?", 1, 1).
		Select("id", "app_key", "app_name", "icon", "redirect_uris").
		Order("created_at ASC").
		Find(&apps).Error

	if err != nil {
		global.Log.Error("获取公开应用列表失败", zap.Error(err))
		return nil, err
	}

	// 判断环境：使用工具方法
	isProduction := core.IsProduction()

	// 转换为响应结构，根据环境选择 redirect_uri
	result := make([]response.PublicApplicationInfo, len(apps))
	for i, app := range apps {
		// 分割 redirect_uris 字符串
		var redirectURI string
		if app.RedirectURIs != "" {
			uris := strings.Split(app.RedirectURIs, ",")
			var validURIs []string
			for _, uri := range uris {
				trimmed := strings.TrimSpace(uri)
				if trimmed != "" {
					validURIs = append(validURIs, trimmed)
				}
			}

			// 根据环境选择 redirect_uri
			// 正式环境：取第一个
			// 测试环境：取第二个（如果存在），否则取第一个
			if len(validURIs) > 0 {
				if isProduction {
					redirectURI = validURIs[0]
				} else {
					if len(validURIs) > 1 {
						redirectURI = validURIs[1]
					} else {
						redirectURI = validURIs[0]
					}
				}
			}
		}

		result[i] = response.PublicApplicationInfo{
			ID:          app.ID,
			AppKey:      app.AppKey,
			AppName:     app.AppName,
			Icon:        app.Icon,
			RedirectURI: redirectURI,
		}
	}

	return result, nil
}
