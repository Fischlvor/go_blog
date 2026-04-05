package admin

import (
	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/shared"
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
	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"avatar":                 a.cfg.Website.Avatar,
		"logo":                   a.cfg.Website.Logo,
		"full_logo":              a.cfg.Website.FullLogo,
		"title":                  a.cfg.Website.Title,
		"slogan":                 a.cfg.Website.Slogan,
		"slogan_en":              a.cfg.Website.SloganEn,
		"description":            a.cfg.Website.Description,
		"version":                a.cfg.Website.Version,
		"created_at":             a.cfg.Website.CreatedAt,
		"icp_filing":             a.cfg.Website.ICPFiling,
		"public_security_filing": a.cfg.Website.PublicSecurityFiling,
		"bilibili_url":           a.cfg.Website.BilibiliURL,
		"github_url":             a.cfg.Website.GithubURL,
		"steam_url":              a.cfg.Website.SteamURL,
		"name":                   a.cfg.Website.Name,
		"job":                    a.cfg.Website.Job,
		"address":                a.cfg.Website.Address,
		"email":                  a.cfg.Website.Email,
		"qq_image":               a.cfg.Website.QQImage,
		"wechat_image":           a.cfg.Website.WechatImage,
	}))
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
