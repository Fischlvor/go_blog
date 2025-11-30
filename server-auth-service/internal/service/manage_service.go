package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	database "auth-service/internal/model/database"
	"auth-service/internal/model/request"
	"auth-service/internal/model/response"
	"auth-service/pkg/global"
)

type ManageService struct{}

// GetUserDevices 获取用户设备列表（包含应用信息）
func (s *ManageService) GetUserDevices(userUUID uuid.UUID, currentDeviceID string) ([]response.DeviceInfo, error) {
	var devices []database.SSODevice
	err := global.DB.Where("user_uuid = ? AND status = 1", userUUID).
		Order("last_active_at DESC").
		Find(&devices).Error

	if err != nil {
		global.Log.Error("查询用户设备失败",
			zap.Error(err),
			zap.String("user_uuid", userUUID.String()),
		)
		return nil, fmt.Errorf("查询设备列表失败: %w", err)
	}

	// 获取所有应用信息
	var apps []database.SSOApplication
	global.DB.Find(&apps)
	appMap := make(map[uint]database.SSOApplication)
	for _, app := range apps {
		appMap[app.ID] = app
	}

	result := make([]response.DeviceInfo, len(devices))
	for i, device := range devices {
		app := appMap[device.AppID]
		result[i] = response.DeviceInfo{
			ID:           device.ID,
			DeviceID:     device.DeviceID,
			DeviceName:   device.DeviceName,
			DeviceType:   device.DeviceType,
			IPAddress:    device.IPAddress,
			LastActiveAt: device.LastActiveAt,
			Status:       device.Status,
			IsCurrent:    device.DeviceID == currentDeviceID, // 标记当前设备
			CreatedAt:    device.CreatedAt,
			AppName:      app.AppName,
			AppKey:       app.AppKey,
		}
	}

	global.Log.Info("获取用户设备列表成功",
		zap.String("user_uuid", userUUID.String()),
		zap.Int("device_count", len(result)),
		zap.String("current_device", currentDeviceID),
	)

	return result, nil
}

// KickDevice 踢出指定设备
func (s *ManageService) KickDevice(c *gin.Context, userUUID uuid.UUID, deviceID string) error {
	// 验证设备归属
	var device database.SSODevice
	err := global.DB.Where("user_uuid = ? AND device_id = ? AND status = 1",
		userUUID, deviceID).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("设备不存在或已离线")
		}
		return fmt.Errorf("查询设备失败: %w", err)
	}

	// 使用现有的踢出设备方法
	authService := &AuthService{}
	err = authService.KickDevice(c, userUUID, deviceID, device.AppID, "manual_kick", "管理员手动踢出")
	if err != nil {
		return fmt.Errorf("踢出设备失败: %w", err)
	}

	global.Log.Info("手动踢出设备成功",
		zap.String("user_uuid", userUUID.String()),
		zap.String("device_id", deviceID),
		zap.String("device_name", device.DeviceName),
	)

	return nil
}

// SSOLogout 退出当前设备SSO（撤销Token和更新设备状态）
func (s *ManageService) SSOLogout(c *gin.Context, userUUID uuid.UUID, deviceID string) error {
	// 验证设备归属
	var device database.SSODevice
	err := global.DB.Where("user_uuid = ? AND device_id = ? AND status = 1",
		userUUID, deviceID).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("设备不存在或已离线")
		}
		return fmt.Errorf("查询设备失败: %w", err)
	}

	// 1. 删除 RefreshToken
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", userUUID.String(), deviceID)
	global.Redis.Del(refreshTokenKey)

	// 2. 更新设备状态
	result := global.DB.Model(&database.SSODevice{}).
		Where("user_uuid = ? AND device_id = ?", userUUID, deviceID).
		Update("status", 0)
	if result.Error != nil {
		return fmt.Errorf("更新设备状态失败: %w", result.Error)
	}

	// 3. 记录SSO退出日志
	authService := &AuthService{}
	authService.LogActionWithContext(c, userUUID, device.AppID, "sso_logout", deviceID, "SSO退出成功", 1)

	global.Log.Info("SSO退出成功",
		zap.String("user_uuid", userUUID.String()),
		zap.String("device_id", deviceID),
		zap.String("device_name", device.DeviceName),
	)

	return nil
}

// AppLogout 退出应用（只撤销Token和更新设备状态，不清除SSO Session）
func (s *ManageService) AppLogout(c *gin.Context, userUUID uuid.UUID, deviceID string) error {
	// 验证设备归属
	var device database.SSODevice
	err := global.DB.Where("user_uuid = ? AND device_id = ? AND status = 1",
		userUUID, deviceID).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("设备不存在或已离线")
		}
		return fmt.Errorf("查询设备失败: %w", err)
	}

	// 1. 删除 RefreshToken（让当前会话失效）
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", userUUID.String(), deviceID)
	global.Redis.Del(refreshTokenKey)

	// 注意：应用退出不删除设备记录，设备仍然存在，只是当前会话结束
	// 只有"踢出设备"或"SSO退出"才需要更新设备状态

	// 2. 记录应用退出日志
	authService := &AuthService{}
	authService.LogActionWithContext(c, userUUID, device.AppID, "logout", deviceID, "应用退出成功", 1)

	global.Log.Info("应用退出成功",
		zap.String("user_uuid", userUUID.String()),
		zap.String("device_id", deviceID),
		zap.String("device_name", device.DeviceName),
	)

	return nil
}

// LogoutAllDevices 退出所有设备
func (s *ManageService) LogoutAllDevices(c *gin.Context, userUUID uuid.UUID) (int, error) {
	// 查询所有活跃设备
	var devices []database.SSODevice
	err := global.DB.Where("user_uuid = ? AND status = 1", userUUID).Find(&devices).Error
	if err != nil {
		return 0, fmt.Errorf("查询用户设备失败: %w", err)
	}

	if len(devices) == 0 {
		return 0, errors.New("没有活跃设备")
	}

	// 批量处理（个位数设备，直接循环）
	kickedCount := 0
	authService := &AuthService{}

	for _, device := range devices {
		err := authService.KickDevice(c, userUUID, device.DeviceID, device.AppID, "logout_all", "退出所有设备")
		if err != nil {
			global.Log.Error("踢出设备失败",
				zap.Error(err),
				zap.String("device_id", device.DeviceID),
			)
		} else {
			kickedCount++
		}
	}

	global.Log.Info("批量退出设备完成",
		zap.String("user_uuid", userUUID.String()),
		zap.Int("total_devices", len(devices)),
		zap.Int("kicked_count", kickedCount),
	)

	return kickedCount, nil
}

// GetOperationLogs 获取操作日志
func (s *ManageService) GetOperationLogs(userUUID uuid.UUID, params request.LogQueryParams) (*response.PageResponse, error) {
	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	// 构建查询条件
	query := global.DB.Model(&database.SSOLoginLog{}).Where("user_uuid = ?", userUUID)

	// 操作类型筛选
	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}

	// 时间范围筛选
	if params.StartTime != "" {
		query = query.Where("created_at >= ?", params.StartTime)
	}
	if params.EndTime != "" {
		query = query.Where("created_at <= ?", params.EndTime)
	}

	// 查询总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("查询日志总数失败: %w", err)
	}

	// 分页查询
	var logs []database.SSOLoginLog
	offset := (params.Page - 1) * params.PageSize
	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("查询操作日志失败: %w", err)
	}

	// 转换为响应格式
	logInfos := make([]response.LogInfo, len(logs))
	for i, log := range logs {
		logInfos[i] = response.LogInfo{
			ID:        log.ID,
			Action:    log.Action,
			DeviceID:  log.DeviceID,
			IPAddress: log.IPAddress,
			Status:    log.Status,
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		}
	}

	return &response.PageResponse{
		List:     logInfos,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}
