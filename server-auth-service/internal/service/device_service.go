package service

import (
	"auth-service/internal/model/database"
	"auth-service/internal/model/response"
	"auth-service/pkg/global"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type DeviceService struct{}

// GetDevices 获取设备列表（基于 UUID）
func (s *DeviceService) GetDevices(userUUID uuid.UUID, appID string, currentDeviceID string) ([]response.DeviceInfo, error) {
	var devices []database.SSODevice

	query := global.DB.Where("user_uuid = ? AND status = 1", userUUID)
	if appID != "" && appID != "all" {
		query = query.Where("app_id = ?", appID)
	}

	if err := query.Order("last_active_at DESC").Find(&devices).Error; err != nil {
		return nil, err
	}

	var result []response.DeviceInfo
	for _, device := range devices {
		result = append(result, response.DeviceInfo{
			DeviceID:     device.DeviceID,
			DeviceName:   device.DeviceName,
			DeviceType:   device.DeviceType,
			IPAddress:    device.IPAddress,
			LastActiveAt: device.LastActiveAt,
			Status:       device.Status,
			IsCurrent:    device.DeviceID == currentDeviceID,
		})
	}

	return result, nil
}

// KickDevice 踢出设备（基于 UUID）
func (s *DeviceService) KickDevice(c *gin.Context, userUUID uuid.UUID, deviceID, appID string, currentDeviceID string) error {
	// 不允许踢出当前设备
	if deviceID == currentDeviceID {
		return errors.New("不能移除当前设备")
	}

	// 验证设备所有权
	var device database.SSODevice
	query := global.DB.Where("device_id = ? AND user_uuid = ?", deviceID, userUUID)
	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}

	if err := query.First(&device).Error; err != nil {
		return errors.New("设备不存在或无权操作")
	}

	// 将该设备的refresh_token从Redis删除（基于 UUID）
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", userUUID.String(), deviceID)
	global.Redis.Del(refreshTokenKey)

	// 将设备标记为黑名单
	deviceBlacklistKey := "device:blacklist:" + deviceID
	global.Redis.Set(deviceBlacklistKey, "*", 0) // 永久黑名单

	// 更新设备状态
	if err := global.DB.Model(&device).Update("status", 0).Error; err != nil {
		return err
	}

	// 记录日志
	authService := &AuthService{}
	authService.LogActionWithContext(c, userUUID, device.AppID, "kick", deviceID, "设备被踢出", 1)

	return nil
}
