package api

import (
	"auth-service/internal/middleware"
	"auth-service/internal/model/response"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type DeviceApi struct {
}

// GetDevices 获取设备列表
func (h *DeviceApi) GetDevices(c *gin.Context) {
	userUUID := middleware.GetUserUUID(c)
	currentDeviceID := middleware.GetDeviceID(c)
	appID := c.Query("app_id") // 可选：指定应用ID

	if userUUID == uuid.Nil {
		response.Unauthorized(c, "未登录")
		return
	}
	devices, err := deviceService.GetDevices(userUUID, appID, currentDeviceID)
	if err != nil {
		response.Error(c, 2001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"devices": devices,
		"total":   len(devices),
	})
}

// KickDevice 踢出设备
func (h *DeviceApi) KickDevice(c *gin.Context) {
	userUUID := middleware.GetUserUUID(c)
	currentDeviceID := middleware.GetDeviceID(c)
	deviceID := c.Param("device_id")
	appID := c.Query("app_id") // 可选：指定应用ID

	if deviceID == "" {
		response.BadRequest(c, "device_id不能为空")
		return
	}

	if userUUID == uuid.Nil {
		response.Unauthorized(c, "未登录")
		return
	}
	if err := deviceService.KickDevice(c, userUUID, deviceID, appID, currentDeviceID); err != nil {
		response.Error(c, 2002, err.Error())
		return
	}

	response.SuccessMsg(c, "设备已移除", nil)
}
