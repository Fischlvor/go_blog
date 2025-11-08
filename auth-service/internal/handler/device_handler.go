package handler

import (
	"auth-service/internal/service"
	"auth-service/pkg/middleware"
	"auth-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService *service.DeviceService
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		deviceService: &service.DeviceService{},
	}
}

// GetDevices 获取设备列表
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	userID := middleware.GetUserID(c)
	currentDeviceID := middleware.GetDeviceID(c)
	appID := c.Query("app_id") // 可选：指定应用ID

	devices, err := h.deviceService.GetDevices(userID, appID, currentDeviceID)
	if err != nil {
		utils.Error(c, 2001, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"devices": devices,
		"total":   len(devices),
	})
}

// KickDevice 踢出设备
func (h *DeviceHandler) KickDevice(c *gin.Context) {
	userID := middleware.GetUserID(c)
	currentDeviceID := middleware.GetDeviceID(c)
	deviceID := c.Param("device_id")
	appID := c.Query("app_id") // 可选：指定应用ID

	if deviceID == "" {
		utils.BadRequest(c, "device_id不能为空")
		return
	}

	if err := h.deviceService.KickDevice(userID, deviceID, appID, currentDeviceID); err != nil {
		utils.Error(c, 2002, err.Error())
		return
	}

	utils.SuccessMsg(c, "设备已移除", nil)
}
