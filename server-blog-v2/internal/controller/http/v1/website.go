package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/controller/http/v1/response"
)

// getWebsiteLogo 获取网站 Logo。
func (v *V1) getWebsiteLogo(c fiber.Ctx) error {
	return shared.WriteSuccess(c, shared.WithData(map[string]string{
		"logo": "/image/logo.svg",
	}))
}

// getWebsiteTitle 获取网站标题。
func (v *V1) getWebsiteTitle(c fiber.Ctx) error {
	info := v.website.GetInfo(c.Context())
	return shared.WriteSuccess(c, shared.WithData(map[string]string{
		"title": info.Title.Value,
	}))
}

// getWebsiteInfo 获取网站信息。
func (v *V1) getWebsiteInfo(c fiber.Ctx) error {
	info := v.website.GetInfo(c.Context())
	return shared.WriteSuccess(c, shared.WithData(info))
}

// getWebsiteCarousel 获取轮播图。
func (v *V1) getWebsiteCarousel(c fiber.Ctx) error {
	carousel, err := v.website.GetCarousel(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - website - getWebsiteCarousel")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorNotImplemented, "failed to get carousel")
	}
	return shared.WriteSuccess(c, shared.WithData(carousel))
}

// getWebsiteNews 获取热搜新闻。
func (v *V1) getWebsiteNews(c fiber.Ctx) error {
	source := c.Query("source", "baidu")
	data, err := v.website.GetNews(c.Context(), source)
	if err != nil {
		v.logger.Error(err, "http - v1 - website - getWebsiteNews")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorNotImplemented, "failed to get news")
	}
	return shared.WriteSuccess(c, shared.WithData(data))
}

// getWebsiteCalendar 获取日历信息。
func (v *V1) getWebsiteCalendar(c fiber.Ctx) error {
	data, err := v.website.GetCalendar(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - website - getWebsiteCalendar")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorNotImplemented, "failed to get calendar")
	}
	return shared.WriteSuccess(c, shared.WithData(data))
}

// getWebsiteFooterLink 获取页脚链接。
func (v *V1) getWebsiteFooterLink(c fiber.Ctx) error {
	links, err := v.website.GetFooterLinks(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - website - getWebsiteFooterLink")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorNotImplemented, "failed to get footer links")
	}
	return shared.WriteSuccess(c, shared.WithData(links))
}

// ==================== 表情相关 ====================

// getEmojiConfig 获取表情配置。
func (v *V1) getEmojiConfig(c fiber.Ctx) error {
	config, err := v.emoji.GetConfig(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - emoji - getEmojiConfig")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorNotImplemented, "failed to get emoji config")
	}
	return shared.WriteSuccess(c, shared.WithData(config))
}

// getEmojiGroups 获取表情组列表。
func (v *V1) getEmojiGroups(c fiber.Ctx) error {
	groups, err := v.emoji.ListGroups(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - emoji - getEmojiGroups")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorNotImplemented, "failed to get emoji groups")
	}
	return shared.WriteSuccess(c, shared.WithData(groups))
}

// ==================== 广告相关 ====================

// getAdvertisementInfo 获取广告信息。
func (v *V1) getAdvertisementInfo(c fiber.Ctx) error {
	info, err := v.advertisement.GetInfo(c.Context())
	if err != nil {
		v.logger.Error(err, "http - v1 - advertisement - getAdvertisementInfo")
		return shared.WriteError(c, http.StatusInternalServerError, response.ErrorNotImplemented, "failed to get advertisement info")
	}
	return shared.WriteSuccess(c, shared.WithData(info))
}
