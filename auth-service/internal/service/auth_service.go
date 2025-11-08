package service

import (
	"auth-service/internal/model/entity"
	"auth-service/internal/model/request"
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

type AuthService struct{}

// getAppByKey 通过app_key获取应用信息
func (s *AuthService) getAppByKey(appKey string) (*entity.SSOApplication, error) {
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

// Register 用户注册
func (s *AuthService) Register(req request.RegisterRequest) error {
	// 检查邮箱是否已注册
	var existUser entity.SSOUser
	err := global.DB.Where("email = ?", req.Email).First(&existUser).Error
	if err == nil {
		return errors.New("邮箱已被注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 哈希密码
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := entity.SSOUser{
		UUID:           uuid.Must(uuid.NewV4()),
		Username:       req.Email, // 默认用邮箱作为用户名
		PasswordHash:   passwordHash,
		Email:          req.Email,
		Nickname:       req.Nickname,
		Avatar:         "/default_avatar.jpg",
		Status:         1,
		RegisterSource: types.Email, // ✅ 使用枚举
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
		return fmt.Errorf("创建用户失败: %w", err)
	}

	// 查询应用ID
	app, err := s.getAppByKey(req.AppID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("应用查询失败: %w", err)
	}

	// 自动授权访问指定应用
	userAppRelation := entity.UserAppRelation{
		UserUUID: user.UUID,
		AppID:    app.ID,
		Status:   1,
	}
	if err := tx.Create(&userAppRelation).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("授权应用失败: %w", err)
	}

	return tx.Commit().Error
}

// Login 用户登录
func (s *AuthService) Login(req request.LoginRequest, ipAddress, userAgent string) (*response.TokenResponse, error) {
	// 查询用户
	var user entity.SSOUser
	err := global.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邮箱或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !crypto.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("邮箱或密码错误")
	}

	// 检查用户状态
	if user.Status == 2 {
		return nil, errors.New("账号已被禁用，请联系管理员")
	}
	if user.Status == 3 {
		return nil, errors.New("账号已注销")
	}

	// 查询应用
	app, err := s.getAppByKey(req.AppID)
	if err != nil {
		return nil, err
	}

	// 检查应用权限
	var userAppRelation entity.UserAppRelation
	err = global.DB.Where("user_uuid = ? AND app_id = ?", user.UUID, app.ID).First(&userAppRelation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("您无权访问此应用")
		}
		return nil, err
	}
	if userAppRelation.Status == 2 {
		return nil, errors.New("您无权访问此应用")
	}

	// 生成设备ID（如果前端没提供）
	deviceID := req.DeviceID
	if deviceID == "" {
		deviceID = uuid.Must(uuid.NewV4()).String()
	}

	// 检查该设备是否已存在（用户+应用+设备的组合唯一）
	var existDevice entity.SSODevice
	err = global.DB.Where("user_id = ? AND app_id = ? AND device_id = ?", user.ID, app.ID, deviceID).First(&existDevice).Error
	isNewDevice := errors.Is(err, gorm.ErrRecordNotFound)

	if isNewDevice {
		// 新设备，检查设备数量
		var deviceCount int64
		global.DB.Model(&entity.SSODevice{}).
			Where("user_id = ? AND app_id = ? AND status = 1", user.ID, app.ID).
			Count(&deviceCount)

		if int(deviceCount) >= app.MaxDevices {
			return nil, fmt.Errorf("设备数已达上限(%d台)，请先在设备管理中移除其他设备", app.MaxDevices)
		}

		// 创建新设备
		device := entity.SSODevice{
			UserID:       user.ID,
			AppID:        app.ID,
			DeviceID:     deviceID,
			DeviceName:   req.DeviceName,
			DeviceType:   req.DeviceType,
			IPAddress:    ipAddress,
			UserAgent:    userAgent,
			LastActiveAt: time.Now(),
			Status:       1,
		}
		if err := global.DB.Create(&device).Error; err != nil {
			return nil, fmt.Errorf("注册设备失败: %w", err)
		}
	} else {
		// 更新现有设备
		global.DB.Model(&existDevice).Updates(map[string]interface{}{
			"device_name":    req.DeviceName,
			"device_type":    req.DeviceType,
			"last_active_at": time.Now(),
			"ip_address":     ipAddress,
			"user_agent":     userAgent,
			"status":         1,
		})
	}

	// 生成Token
	accessTokenDuration, _ := utils.ParseDuration(global.Config.JWT.AccessTokenExpiryTime)
	refreshTokenDuration, _ := utils.ParseDuration(global.Config.JWT.RefreshTokenExpiryTime)

	accessToken, err := jwt.CreateAccessToken(
		user.ID,
		user.UUID,
		req.AppID,
		deviceID,
		accessTokenDuration,
		global.Config.JWT.Issuer,
		global.RSAPrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("生成AccessToken失败: %w", err)
	}

	refreshToken, err := jwt.CreateRefreshToken(
		user.ID,
		req.AppID,
		deviceID,
		refreshTokenDuration,
		global.Config.JWT.Issuer,
		global.RSAPrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("生成RefreshToken失败: %w", err)
	}

	// 将RefreshToken存储到Redis（用于撤销检查）
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
		Message:   "登录成功",
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

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(req request.RefreshTokenRequest) (*response.TokenResponse, error) {
	// 验证client_secret
	app, err := s.getAppByKey(req.ClientID)
	if err != nil {
		return nil, err
	}
	if app.AppSecret != req.ClientSecret {
		return nil, errors.New("client_secret错误")
	}

	// 解析RefreshToken
	claims, err := jwt.ParseRefreshToken(req.RefreshToken, global.RSAPublicKey)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, errors.New("refresh_token已过期，请重新登录")
		}
		return nil, errors.New("refresh_token无效")
	}

	// 检查Redis中是否存在（未被撤销）
	refreshTokenKey := fmt.Sprintf("refresh_token:%d:%s", claims.UserID, claims.DeviceID)
	storedToken, err := global.Redis.Get(refreshTokenKey).Result()
	if err != nil || storedToken != req.RefreshToken {
		return nil, errors.New("refresh_token已被撤销，请重新登录")
	}

	// 检查用户状态
	var user entity.SSOUser
	if err := global.DB.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用或注销")
	}

	// 检查设备状态
	var device entity.SSODevice
	err = global.DB.Where("user_id = ? AND device_id = ?", claims.UserID, claims.DeviceID).First(&device).Error
	if err != nil || device.Status != 1 {
		return nil, errors.New("设备已被移除")
	}

	// 生成新的AccessToken
	accessTokenDuration, _ := utils.ParseDuration(global.Config.JWT.AccessTokenExpiryTime)
	accessToken, err := jwt.CreateAccessToken(
		user.ID,
		user.UUID,
		claims.AppID,
		claims.DeviceID,
		accessTokenDuration,
		global.Config.JWT.Issuer,
		global.RSAPrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("生成AccessToken失败: %w", err)
	}

	// Token轮换：生成新的RefreshToken（可选，更安全）
	refreshTokenDuration, _ := utils.ParseDuration(global.Config.JWT.RefreshTokenExpiryTime)
	newRefreshToken, err := jwt.CreateRefreshToken(
		user.ID,
		claims.AppID,
		claims.DeviceID,
		refreshTokenDuration,
		global.Config.JWT.Issuer,
		global.RSAPrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("生成RefreshToken失败: %w", err)
	}

	// 更新Redis中的RefreshToken
	global.Redis.Set(refreshTokenKey, newRefreshToken, refreshTokenDuration)

	// 更新设备活跃时间
	global.DB.Model(&device).Update("last_active_at", time.Now())

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessTokenDuration.Seconds()),
	}, nil
}

// Logout 登出
func (s *AuthService) Logout(accessToken string) error {
	// 解析Token
	claims, err := jwt.ParseAccessTokenIgnoreExpiry(accessToken, global.RSAPublicKey)
	if err != nil {
		return errors.New("Token无效")
	}

	// 将AccessToken加入黑名单
	blacklistKey := "token:blacklist:" + accessToken
	// 过期时间设为Token剩余有效期
	expiry := time.Until(time.Unix(claims.ExpiresAt.Unix(), 0))
	if expiry > 0 {
		global.Redis.Set(blacklistKey, "1", expiry)
	}

	// 删除RefreshToken
	refreshTokenKey := fmt.Sprintf("refresh_token:%d:%s", claims.UserID, claims.DeviceID)
	global.Redis.Del(refreshTokenKey)

	// 更新设备状态
	global.DB.Model(&entity.SSODevice{}).
		Where("user_id = ? AND device_id = ?", claims.UserID, claims.DeviceID).
		Update("status", 0)

	// 查询应用ID（用于日志）
	var logAppID uint
	if app, err := s.getAppByKey(claims.AppID); err == nil {
		logAppID = app.ID
	}

	// 记录登出日志
	loginLog := entity.SSOLoginLog{
		UserID:   claims.UserID,
		AppID:    logAppID,
		Action:   "logout",
		DeviceID: claims.DeviceID,
		Status:   1,
		Message:  "登出成功",
	}
	global.DB.Create(&loginLog)

	return nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID uint) (*response.UserInfo, error) {
	var user entity.SSOUser
	if err := global.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &response.UserInfo{
		UserID:    user.ID,
		UUID:      user.UUID.String(),
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *AuthService) UpdateUserInfo(userID uint, req request.UpdateUserInfoRequest) error {
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Signature != "" {
		updates["signature"] = req.Signature
	}

	if len(updates) == 0 {
		return errors.New("没有需要更新的内容")
	}

	return global.DB.Model(&entity.SSOUser{}).Where("id = ?", userID).Updates(updates).Error
}

// UpdatePassword 修改密码
func (s *AuthService) UpdatePassword(userID uint, req request.UpdatePasswordRequest) error {
	var user entity.SSOUser
	if err := global.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !crypto.CheckPassword(req.OldPassword, user.PasswordHash) {
		return errors.New("原密码错误")
	}

	// 哈希新密码
	newPasswordHash, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	return global.DB.Model(&user).Update("password_hash", newPasswordHash).Error
}
