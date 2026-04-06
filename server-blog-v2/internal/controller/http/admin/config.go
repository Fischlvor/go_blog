package admin

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/admin/request"
	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/usecase/input"
)

// maskSecret 隐藏敏感信息，只显示前4位。
func maskSecret(s string) string {
	if len(s) <= 4 {
		return "******"
	}
	return s[:4] + "******"
}

// getWebsiteConfig 获取网站配置。
// @Summary 获取网站配置（管理端）
// @Tags Admin.Config
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/config/website [get]
func (a *Admin) getWebsiteConfig(c fiber.Ctx) error {
	info := a.website.GetInfo(c.Context())
	return shared.WriteSuccess(c, shared.WithData(info))
}

// updateWebsiteConfig 更新网站配置。
// @Summary 更新网站配置（管理端）
// @Tags Admin.Config
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.WebsiteConfig true "网站配置"
// @Success 200 {object} shared.Envelope
// @Router /admin/config/website [put]
func (a *Admin) updateWebsiteConfig(c fiber.Ctx) error {
	var req request.WebsiteConfig
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	settings := []input.UpsertSiteSetting{
		toUpsertSiteSetting(req.Avatar, "profile.avatar", "string", "头像"),
		toUpsertSiteSetting(req.Title, "website.title", "string", "网站标题"),
		toUpsertSiteSetting(req.Description, "website.description", "string", "网站描述"),
		toUpsertSiteSetting(req.ProfileIntro, "profile.intro", "string", "个人介绍"),
		toUpsertSiteSetting(req.TechStack, "profile.tech_stack", "json", "技术栈"),
		toUpsertSiteSetting(req.WorkExperiences, "profile.work_experiences", "json", "工作经历"),
		toUpsertSiteSetting(req.Version, "website.version", "string", "网站版本"),
		toUpsertSiteSetting(req.CreatedAt, "website.created_at", "string", "网站创建日期"),
		toUpsertSiteSetting(req.ICPFiling, "website.icp_filing", "string", "ICP备案号"),
		toUpsertSiteSetting(req.BilibiliURL, "profile.bilibili_url", "string", "Bilibili 链接"),
		toUpsertSiteSetting(req.GithubURL, "profile.github_url", "string", "GitHub 链接"),
		toUpsertSiteSetting(req.SteamURL, "profile.steam_url", "string", "Steam 链接"),
		toUpsertSiteSetting(req.Name, "profile.name", "string", "名称"),
		toUpsertSiteSetting(req.Job, "profile.job", "string", "职业"),
		toUpsertSiteSetting(req.Address, "profile.address", "string", "地址"),
		toUpsertSiteSetting(req.Email, "profile.email", "string", "联系邮箱"),
	}

	if err := a.setting.UpdateSiteSettings(c.Context(), settings); err != nil {
		a.logger.Error(err, "http - admin - config - updateWebsiteConfig")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to update website config")
	}

	return shared.WriteSuccess(c)
}

func toUpsertSiteSetting(field request.SettingField, fallbackKey, settingType, description string) input.UpsertSiteSetting {
	settingKey := field.SettingKey
	if settingKey == "" {
		settingKey = fallbackKey
	}
	return input.UpsertSiteSetting{
		SettingKey:   settingKey,
		SettingValue: field.Value,
		SettingType:  settingType,
		Description:  description,
		IsPublic:     true,
	}
}

// getSystemConfig 获取系统配置。
// @Summary 获取系统配置（管理端）
// @Tags Admin.Config
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/config/system [get]
func (a *Admin) getSystemConfig(c fiber.Ctx) error {
	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"use_multipoint":  a.cfg.System.UseMultipoint,
		"sessions_secret": maskSecret(a.cfg.System.SessionsSecret),
		"oss_type":        a.cfg.System.OssType,
	}))
}

// getEmailConfig 获取邮箱配置。
// @Summary 获取邮箱配置（管理端）
// @Tags Admin.Config
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/config/email [get]
func (a *Admin) getEmailConfig(c fiber.Ctx) error {
	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"host":     a.cfg.Email.Host,
		"port":     a.cfg.Email.Port,
		"from":     a.cfg.Email.From,
		"nickname": a.cfg.Email.Nickname,
		"secret":   maskSecret(a.cfg.Email.Secret),
		"is_ssl":   a.cfg.Email.IsSSL,
	}))
}

// getQQConfig 获取QQ登录配置。
// @Summary 获取QQ登录配置（管理端）
// @Tags Admin.Config
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/config/qq [get]
func (a *Admin) getQQConfig(c fiber.Ctx) error {
	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"enable":       a.cfg.QQ.Enable,
		"app_id":       a.cfg.QQ.AppID,
		"app_key":      maskSecret(a.cfg.QQ.AppKey),
		"redirect_uri": a.cfg.QQ.RedirectURI,
	}))
}

// getQiniuConfig 获取七牛云配置。
// @Summary 获取七牛云配置（管理端）
// @Tags Admin.Config
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/config/qiniu [get]
func (a *Admin) getQiniuConfig(c fiber.Ctx) error {
	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"zone":            a.cfg.Qiniu.Zone,
		"bucket":          a.cfg.Qiniu.Bucket,
		"img_path":        a.cfg.Qiniu.Domain,
		"access_key":      maskSecret(a.cfg.Qiniu.AccessKey),
		"secret_key":      maskSecret(a.cfg.Qiniu.SecretKey),
		"use_https":       a.cfg.Qiniu.UseHTTPS,
		"use_cdn_domains": a.cfg.Qiniu.UseCDNDomains,
	}))
}

// getJwtConfig 获取JWT配置。
// @Summary 获取JWT配置（管理端）
// @Tags Admin.Config
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/config/jwt [get]
func (a *Admin) getJwtConfig(c fiber.Ctx) error {
	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"access_token_secret":       maskSecret(a.cfg.Jwt.AccessTokenSecret),
		"refresh_token_secret":      maskSecret(a.cfg.Jwt.RefreshTokenSecret),
		"access_token_expiry_time":  a.cfg.Jwt.AccessTokenExpiryTime,
		"refresh_token_expiry_time": a.cfg.Jwt.RefreshTokenExpiryTime,
		"issuer":                    a.cfg.Jwt.Issuer,
	}))
}

// getGaodeConfig 获取高德配置。
// @Summary 获取高德配置（管理端）
// @Tags Admin.Config
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /admin/config/gaode [get]
func (a *Admin) getGaodeConfig(c fiber.Ctx) error {
	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"enable": a.cfg.Gaode.Enable,
		"key":    maskSecret(a.cfg.Gaode.Key),
	}))
}
