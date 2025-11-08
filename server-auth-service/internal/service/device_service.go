package service

import (
	"auth-service/internal/model/entity"
	"auth-service/internal/model/response"
	"auth-service/pkg/global"
	"errors"
	"fmt"
)

type DeviceService struct{}

// GetDevices 获取设备列表
func (s *DeviceService) GetDevices(userID uint, appID string, currentDeviceID string) ([]response.DeviceInfo, error) {
	var devices []entity.SSODevice

	query := global.DB.Where("user_id = ? AND status = 1", userID)
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

// KickDevice 踢出设备
func (s *DeviceService) KickDevice(userID uint, deviceID, appID string, currentDeviceID string) error {
	// 不允许踢出当前设备
	if deviceID == currentDeviceID {
		return errors.New("不能移除当前设备")
	}

	// 验证设备所有权
	var device entity.SSODevice
	query := global.DB.Where("device_id = ? AND user_id = ?", deviceID, userID)
	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}

	if err := query.First(&device).Error; err != nil {
		return errors.New("设备不存在或无权操作")
	}

	// 将该设备的refresh_token从Redis删除
	refreshTokenKey := fmt.Sprintf("refresh_token:%d:%s", userID, deviceID)
	global.Redis.Del(refreshTokenKey)

	// 将设备标记为黑名单
	deviceBlacklistKey := "device:blacklist:" + deviceID
	global.Redis.Set(deviceBlacklistKey, "*", 0) // 永久黑名单

	// 更新设备状态
	if err := global.DB.Model(&device).Update("status", 0).Error; err != nil {
		return err
	}

	// 记录日志
	loginLog := entity.SSOLoginLog{
		UserID:   userID,
		AppID:    device.AppID,
		Action:   "kick",
		DeviceID: deviceID,
		Status:   1,
		Message:  "设备被踢出",
	}
	global.DB.Create(&loginLog)

	return nil
}
