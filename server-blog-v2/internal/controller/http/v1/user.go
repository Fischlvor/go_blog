package v1

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"server-blog-v2/internal/controller/http/bizcode"
	"server-blog-v2/internal/controller/http/middleware"
	"server-blog-v2/internal/controller/http/shared"
	"server-blog-v2/internal/controller/http/v1/response"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/pkg/gaode"
)

// getProfile 获取用户资料。
func (v *V1) getProfile(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	profile, err := v.user.GetProfile(c.Context(), userUUID)
	if err != nil {
		v.logger.Error(err, "http - v1 - user - getProfile")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to get profile")
	}

	return shared.WriteSuccess(c, shared.WithData(profile))
}

// updateProfile 更新用户资料。
func (v *V1) updateProfile(c fiber.Ctx) error {
	userUUID := middleware.GetUserUUID(c)
	if userUUID == "" {
		return shared.WriteError(c, http.StatusUnauthorized, bizcode.ErrorLoginRequired, "login required")
	}

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteError(c, http.StatusBadRequest, bizcode.ErrorParam, "invalid request body")
	}

	if err := v.user.UpdateProfile(c.Context(), userUUID, input.UpdateProfile{
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}); err != nil {
		v.logger.Error(err, "http - v1 - user - updateProfile")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to update profile")
	}

	return shared.WriteSuccess(c, shared.WithMsg("profile updated"))
}

// getUserCard 获取用户名片（公开）。
// @Summary 获取用户名片
// @Tags V1.User
// @Produce json
// @Param uuid query string true "用户 UUID"
// @Success 200 {object} shared.Envelope
// @Router /v1/user/card [get]
func (v *V1) getUserCard(c fiber.Ctx) error {
	uuid := c.Query("uuid")
	if uuid == "" {
		return shared.WriteError(c, http.StatusBadRequest, response.ErrorParam, "uuid is required")
	}

	profile, err := v.user.GetProfile(c.Context(), uuid)
	if err != nil {
		v.logger.Error(err, "http - v1 - user - getUserCard")
		return shared.WriteError(c, http.StatusNotFound, response.ErrorUserNotFound, "user not found")
	}

	return shared.WriteSuccess(c, shared.WithData(map[string]interface{}{
		"uuid":      profile.UUID,
		"username":  profile.Nickname,
		"avatar":    profile.Avatar,
		"address":   "",
		"signature": profile.Signature,
	}))
}

// getUserWeather 获取用户天气。
// @Summary 获取用户天气
// @Tags V1.User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /v1/user/weather [get]
func (v *V1) getUserWeather(c fiber.Ctx) error {
	ip := c.IP()

	// 如果没有配置高德 key，返回默认值
	if v.cfg.Gaode.Key == "" {
		return shared.WriteSuccess(c, shared.WithData("天气服务未配置"))
	}

	client := gaode.NewClient(v.cfg.Gaode.Key)
	weather, err := client.GetWeatherString(ip)
	if err != nil {
		v.logger.Error(err, "http - v1 - user - getUserWeather")
		return shared.WriteSuccess(c, shared.WithData("获取天气失败"))
	}

	return shared.WriteSuccess(c, shared.WithData(weather))
}

// getUserChart 获取用户图表数据。
// @Summary 获取用户活动图表
// @Tags V1.User
// @Security BearerAuth
// @Produce json
// @Param date query int false "天数" default(7)
// @Success 200 {object} shared.Envelope
// @Router /v1/user/chart [get]
func (v *V1) getUserChart(c fiber.Ctx) error {
	days := 7
	if d := c.Query("date"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 30 {
			days = parsed
		}
	}

	result, err := v.user.GetChart(c.Context(), days)
	if err != nil {
		v.logger.Error(err, "http - v1 - user - getUserChart")
		return shared.WriteError(c, http.StatusInternalServerError, bizcode.ErrorDatabase, "failed to get chart data")
	}

	return shared.WriteSuccess(c, shared.WithData(result))
}

// logout 用户登出。
// @Summary 用户登出
// @Tags V1.User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} shared.Envelope
// @Router /v1/user/logout [post]
func (v *V1) logout(c fiber.Ctx) error {
	// SSO 模式下，登出由前端清除 token 完成
	// 后端只返回成功响应
	return shared.WriteSuccess(c, shared.WithMsg("logout success"))
}
