package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"auth-service/internal/service"
	"auth-service/pkg/global"
	"auth-service/pkg/jwt"
	"auth-service/pkg/middleware"
	"auth-service/pkg/utils"
)

type ManageHandler struct {
	manageService *service.ManageService
}

func NewManageHandler() *ManageHandler {
	return &ManageHandler{
		manageService: &service.ManageService{},
	}
}

// GetDevices 获取用户设备列表
func (h *ManageHandler) GetDevices(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		utils.Error(c, 1001, "用户未登录")
		return
	}

	// 从Token获取当前设备ID
	token := c.GetHeader("Authorization")
	if len(token) <= 7 || token[:7] != "Bearer " {
		utils.Error(c, 1001, "Token格式错误")
		return
	}
	token = token[7:]

	claims, err := jwt.ParseAccessToken(token, global.RSAPublicKey)
	if err != nil {
		utils.Error(c, 1001, "Token解析失败")
		return
	}

	// 获取设备列表
	devices, err := h.manageService.GetUserDevices(userUUID, claims.DeviceID)
	if err != nil {
		global.Log.Error("获取设备列表失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
		)
		utils.Error(c, 2004, "获取设备列表失败")
		return
	}

	utils.Success(c, devices)
}

// KickDevice 踢出指定设备
func (h *ManageHandler) KickDevice(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		utils.Error(c, 1001, "用户未登录")
		return
	}

	// 获取请求参数
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查是否尝试踢出自己
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		claims, err := jwt.ParseAccessToken(token[7:], global.RSAPublicKey)
		if err == nil && claims.DeviceID == req.DeviceID {
			utils.Error(c, 2002, "不能踢出当前设备")
			return
		}
	}

	// 踢出设备
	err := h.manageService.KickDevice(c, userUUID, req.DeviceID)
	if err != nil {
		global.Log.Error("踢出设备失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
			zap.String("device_id", req.DeviceID),
		)
		utils.Error(c, 2004, err.Error())
		return
	}

	utils.Success(c, "设备已踢出")
}

// SSOLogout SSO退出（清除Session）
func (h *ManageHandler) SSOLogout(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		utils.Error(c, 1001, "用户未登录")
		return
	}

	// 从Token获取设备ID
	token := c.GetHeader("Authorization")
	if len(token) <= 7 || token[:7] != "Bearer " {
		utils.Error(c, 1001, "Token格式错误")
		return
	}

	claims, err := jwt.ParseAccessToken(token[7:], global.RSAPublicKey)
	if err != nil {
		utils.Error(c, 1001, "Token解析失败")
		return
	}

	// 撤销Token和更新设备状态
	err = h.manageService.SSOLogout(c, userUUID, claims.DeviceID)
	if err != nil {
		global.Log.Error("SSO退出失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
			zap.String("device_id", claims.DeviceID),
		)
		utils.Error(c, 2004, err.Error())
		return
	}

	// 清除当前Session
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		global.Log.Error("清除Session失败", zap.Error(err))
	}

	global.Log.Info("SSO退出成功",
		zap.String("user_uuid", userUUID.String()),
		zap.String("device_id", claims.DeviceID),
	)

	utils.Success(c, "SSO退出成功")
}

// LogoutAllDevices 退出所有设备
func (h *ManageHandler) LogoutAllDevices(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		utils.Error(c, 1001, "用户未登录")
		return
	}

	// 退出所有设备
	kickedCount, err := h.manageService.LogoutAllDevices(c, userUUID)
	if err != nil {
		global.Log.Error("退出所有设备失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
		)
		utils.Error(c, 2004, err.Error())
		return
	}

	// 清除当前Session
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		global.Log.Error("清除Session失败", zap.Error(err))
	}

	global.Log.Info("退出所有设备成功",
		zap.String("user_uuid", userUUID.String()),
		zap.Int("kicked_count", kickedCount),
	)

	utils.Success(c, gin.H{
		"message":      "已退出所有设备",
		"kicked_count": kickedCount,
	})
}

// GetLogs 获取操作日志
func (h *ManageHandler) GetLogs(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		utils.Error(c, 1001, "用户未登录")
		return
	}

	// 绑定查询参数
	var params service.LogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取操作日志
	result, err := h.manageService.GetOperationLogs(userUUID, params)
	if err != nil {
		global.Log.Error("获取操作日志失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
		)
		utils.Error(c, 2004, "获取操作日志失败")
		return
	}

	utils.Success(c, result)
}

// GetProfile 获取当前用户信息
func (h *ManageHandler) GetProfile(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		utils.Error(c, 1001, "用户未登录")
		return
	}

	// 获取用户信息（复用现有的AuthService方法）
	authService := &service.AuthService{}
	userInfo, err := authService.GetUserInfoByUUID(userUUID)
	if err != nil {
		global.Log.Error("获取用户信息失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
		)
		utils.Error(c, 2004, "获取用户信息失败")
		return
	}

	utils.Success(c, userInfo)
}
