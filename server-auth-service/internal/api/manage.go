package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"auth-service/internal/middleware"
	"auth-service/internal/model/request"
	"auth-service/internal/model/response"
	"auth-service/internal/service"
	"auth-service/pkg/global"
	"auth-service/pkg/jwt"
)

type ManageApi struct {
}

// GetDevices 获取用户设备列表
func (h *ManageApi) GetDevices(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Error(c, 1001, "用户未登录")
		return
	}

	// 从Token获取当前设备ID
	token := c.GetHeader("Authorization")
	if len(token) <= 7 || token[:7] != "Bearer " {
		response.Error(c, 1001, "Token格式错误")
		return
	}
	token = token[7:]

	claims, err := jwt.ParseAccessToken(token, global.RSAPublicKey)
	if err != nil {
		response.Error(c, 1001, "Token解析失败")
		return
	}

	// 获取设备列表
	devices, err := manageService.GetUserDevices(userUUID, claims.DeviceID)
	if err != nil {
		global.Log.Error("获取设备列表失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
		)
		response.Error(c, 2004, "获取设备列表失败")
		return
	}

	response.Success(c, devices)
}

// KickDevice 踢出指定设备
func (h *ManageApi) KickDevice(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Error(c, 1001, "用户未登录")
		return
	}

	// 获取请求参数
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查是否尝试踢出自己
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		claims, err := jwt.ParseAccessToken(token[7:], global.RSAPublicKey)
		if err == nil && claims.DeviceID == req.DeviceID {
			response.Error(c, 2002, "不能踢出当前设备")
			return
		}
	}

	// 踢出设备
	err := manageService.KickDevice(c, userUUID, req.DeviceID)
	if err != nil {
		global.Log.Error("踢出设备失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
			zap.String("device_id", req.DeviceID),
		)
		response.Error(c, 2004, err.Error())
		return
	}

	response.Success(c, "设备已踢出")
}

// Logout 退出manage应用（只退出当前应用，不影响SSO）
func (h *ManageApi) Logout(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Error(c, 1001, "用户未登录")
		return
	}

	// 从Token获取设备ID
	token := c.GetHeader("Authorization")
	if len(token) <= 7 || token[:7] != "Bearer " {
		response.Error(c, 1001, "Token格式错误")
		return
	}

	claims, err := jwt.ParseAccessToken(token[7:], global.RSAPublicKey)
	if err != nil {
		response.Error(c, 1001, "Token解析失败")
		return
	}

	// 只撤销manage应用的Token和更新设备状态，不清除SSO Session
	err = manageService.AppLogout(c, userUUID, claims.DeviceID)
	if err != nil {
		global.Log.Error("manage应用退出失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
			zap.String("device_id", claims.DeviceID),
		)
		response.Error(c, 2004, err.Error())
		return
	}

	// 注意：不清除SSO Session，保持其他应用的登录状态
	global.Log.Info("manage应用退出成功",
		zap.String("user_uuid", userUUID.String()),
		zap.String("device_id", claims.DeviceID),
	)

	// 返回重定向信息
	response.SuccessMsg(c, "manage退出成功", gin.H{
		"redirect_to": "/login?app_id=manage",
	})
}

// SSOLogout SSO全局退出（清除Session，影响所有应用）
func (h *ManageApi) SSOLogout(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Error(c, 1001, "用户未登录")
		return
	}

	// 从Token获取设备ID
	token := c.GetHeader("Authorization")
	if len(token) <= 7 || token[:7] != "Bearer " {
		response.Error(c, 1001, "Token格式错误")
		return
	}

	claims, err := jwt.ParseAccessToken(token[7:], global.RSAPublicKey)
	if err != nil {
		response.Error(c, 1001, "Token解析失败")
		return
	}

	// 撤销Token和更新设备状态
	err = manageService.SSOLogout(c, userUUID, claims.DeviceID)
	if err != nil {
		global.Log.Error("SSO退出失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
			zap.String("device_id", claims.DeviceID),
		)
		response.Error(c, 2004, err.Error())
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

	// 返回重定向信息，前端根据这个跳转到登录页面
	response.SuccessMsg(c, "SSO退出成功", gin.H{
		"redirect_to": "/login?app_id=manage", // 指定 app_id=manage
	})
}

// LogoutAllDevices 退出所有设备
func (h *ManageApi) LogoutAllDevices(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Error(c, 1001, "用户未登录")
		return
	}

	// 退出所有设备
	kickedCount, err := manageService.LogoutAllDevices(c, userUUID)
	if err != nil {
		global.Log.Error("退出所有设备失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
		)
		response.Error(c, 2004, err.Error())
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

	response.Success(c, gin.H{
		"message":      "已退出所有设备",
		"kicked_count": kickedCount,
	})
}

// GetLogs 获取操作日志
func (h *ManageApi) GetLogs(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Error(c, 1001, "用户未登录")
		return
	}

	// 绑定查询参数
	var params request.LogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取操作日志
	result, err := manageService.GetOperationLogs(userUUID, params)
	if err != nil {
		global.Log.Error("获取操作日志失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
		)
		response.Error(c, 2004, "获取操作日志失败")
		return
	}

	response.Success(c, result)
}

// GetProfile 获取当前用户信息
func (h *ManageApi) GetProfile(c *gin.Context) {
	// 从中间件获取用户UUID
	userUUID := middleware.GetUserUUID(c)
	if userUUID == uuid.Nil {
		response.Error(c, 1001, "用户未登录")
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
		response.Error(c, 2004, "获取用户信息失败")
		return
	}

	response.Success(c, userInfo)
}
