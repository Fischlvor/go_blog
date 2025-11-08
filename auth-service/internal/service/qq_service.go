package service

import (
	"auth-service/internal/model/entity"
	"auth-service/internal/model/response"
	"auth-service/internal/model/types"
	"auth-service/pkg/crypto"
	"auth-service/pkg/global"
	"auth-service/pkg/jwt"
	"auth-service/pkg/utils"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type QQService struct{}

// getAppByKey 通过app_key获取应用信息
func (s *QQService) getAppByKey(appKey string) (*entity.SSOApplication, error) {
	var app entity.SSOApplication
	err := global.DB.Where("app_key = ? AND status = 1", appKey).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("应用不存在或未启用")
		}
		return nil, err
	}
	return &app, nil
}

// QQLogin QQ登录
func (s *QQService) QQLogin(openID, appID, deviceID, deviceName, deviceType, ipAddress, userAgent string, qqUserInfo map[string]interface{}) (*response.TokenResponse, error) {
	// 查找OAuth绑定
	var oauthBinding entity.SSOOAuthBinding
	err := global.DB.Where("provider = ? AND open_id = ?", "qq", openID).First(&oauthBinding).Error

	var user entity.SSOUser
	var isNewUser bool

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 新用户，创建账号
		isNewUser = true

		// 从QQ获取昵称和头像
		nickname := "QQ用户"
		avatar := "/default_avatar.jpg"
		if name, ok := qqUserInfo["nickname"].(string); ok && name != "" {
			nickname = name
		}
		if figureurl, ok := qqUserInfo["figureurl_qq_2"].(string); ok && figureurl != "" {
			avatar = figureurl
		}

		user = entity.SSOUser{
			UUID:           uuid.Must(uuid.NewV4()),
			Username:       fmt.Sprintf("qq_%s", openID),
			PasswordHash:   crypto.HashPasswordDefault(), // 随机密码
			Nickname:       nickname,
			Avatar:         avatar,
			Status:         1,
			RegisterSource: types.QQ, // ✅ 使用枚举
			IsSuperAdmin:   false,
		}

		// 开启事务
		tx := global.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 创建用户
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}

		// 创建OAuth绑定
		oauthBinding = entity.SSOOAuthBinding{
			UserID:   user.ID,
			Provider: "qq",
			OpenID:   openID,
		}
		if err := tx.Create(&oauthBinding).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("创建OAuth绑定失败: %w", err)
		}

		// 查询应用ID
		app, err := s.getAppByKey(appID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("应用查询失败: %w", err)
		}

		// 自动授权访问指定应用
		userAppRelation := entity.UserAppRelation{
			UserUUID: user.UUID,
			AppID:    app.ID,
			Status:   1,
		}
		if err := tx.Create(&userAppRelation).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("授权应用失败: %w", err)
		}

		if err := tx.Commit().Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// 老用户，查询用户信息
		if err := global.DB.First(&user, oauthBinding.UserID).Error; err != nil {
			return nil, errors.New("用户不存在")
		}
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用或注销")
	}

	// 查询应用
	app, err := s.getAppByKey(appID)
	if err != nil {
		return nil, err
	}

	// 检查应用权限
	var userAppRelation entity.UserAppRelation
	err = global.DB.Where("user_uuid = ? AND app_id = ?", user.UUID, app.ID).First(&userAppRelation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 自动授权
			userAppRelation = entity.UserAppRelation{
				UserUUID: user.UUID,
				AppID:    app.ID,
				Status:   1,
			}
			global.DB.Create(&userAppRelation)
		} else {
			return nil, err
		}
	}
	if userAppRelation.Status == 2 {
		return nil, errors.New("您无权访问此应用")
	}

	if deviceID == "" {
		deviceID = uuid.Must(uuid.NewV4()).String()
	}

	// 检查该设备是否已存在（用户+应用+设备的组合唯一）
	var existDevice entity.SSODevice
	err = global.DB.Where("user_id = ? AND app_id = ? AND device_id = ?", user.ID, app.ID, deviceID).First(&existDevice).Error
	isNewDevice := errors.Is(err, gorm.ErrRecordNotFound)

	if isNewDevice {
		var deviceCount int64
		global.DB.Model(&entity.SSODevice{}).
			Where("user_id = ? AND app_id = ? AND status = 1", user.ID, app.ID).
			Count(&deviceCount)

		if int(deviceCount) >= app.MaxDevices {
			return nil, fmt.Errorf("设备数已达上限(%d台)", app.MaxDevices)
		}
	}

	device := entity.SSODevice{
		UserID:       user.ID,
		AppID:        app.ID,
		DeviceID:     deviceID,
		DeviceName:   deviceName,
		DeviceType:   deviceType,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		LastActiveAt: time.Now(),
		Status:       1,
	}

	if isNewDevice {
		global.DB.Create(&device)
	} else {
		global.DB.Model(&existDevice).Updates(map[string]interface{}{
			"last_active_at": time.Now(),
			"ip_address":     ipAddress,
			"user_agent":     userAgent,
			"status":         1,
		})
	}

	// 生成Token
	accessTokenDuration, _ := utils.ParseDuration(global.Config.JWT.AccessTokenExpiryTime)
	refreshTokenDuration, _ := utils.ParseDuration(global.Config.JWT.RefreshTokenExpiryTime)

	accessToken, _ := jwt.CreateAccessToken(
		user.ID,
		user.UUID,
		appID,
		deviceID,
		accessTokenDuration,
		global.Config.JWT.Issuer,
		global.RSAPrivateKey,
	)

	refreshToken, _ := jwt.CreateRefreshToken(
		user.ID,
		appID,
		deviceID,
		refreshTokenDuration,
		global.Config.JWT.Issuer,
		global.RSAPrivateKey,
	)

	// 存储RefreshToken
	refreshTokenKey := fmt.Sprintf("refresh_token:%d:%s", user.ID, deviceID)
	global.Redis.Set(refreshTokenKey, refreshToken, refreshTokenDuration)

	// 记录登录日志
	loginLog := entity.SSOLoginLog{
		UserID:    user.ID,
		AppID:     app.ID,
		Action:    "login",
		DeviceID:  deviceID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    1,
		Message:   fmt.Sprintf("QQ登录%s", map[bool]string{true: "（新用户）", false: ""}[isNewUser]),
	}
	global.DB.Create(&loginLog)

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessTokenDuration.Seconds()),
		UserInfo: &response.UserInfo{
			UserID:    user.ID,
			UUID:      user.UUID.String(),
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Email:     user.Email,
			Address:   user.Address,
			Signature: user.Signature,
		},
	}, nil
}
