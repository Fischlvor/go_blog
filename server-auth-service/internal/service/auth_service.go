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
	"go.uber.org/zap"
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
		Avatar:         "https://image.hsk423.cn/blog/aaca0f5eb4d2d98a6ce6dffa99f8254b-20251108151238.jpg",
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

	// 根据登录方式验证
	if req.Password != "" {
		// 密码登录：验证密码
		if !crypto.CheckPassword(req.Password, user.PasswordHash) {
			return nil, errors.New("邮箱或密码错误")
		}
	} else if req.VerificationCode != "" {
		// 验证码登录：验证邮箱验证码
		key := fmt.Sprintf("email_verification_code:%s", req.Email)
		storedCode, err := global.Redis.Get(key).Result()
		if err != nil {
			return nil, errors.New("验证码已过期或不存在")
		}
		if storedCode != req.VerificationCode {
			return nil, errors.New("验证码错误")
		}
		// 验证成功后删除验证码（一次性使用）
		global.Redis.Del(key)
	} else {
		return nil, errors.New("请提供密码或邮箱验证码")
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
	err = global.DB.Where("user_uuid = ? AND app_id = ? AND device_id = ?", user.UUID, app.ID, deviceID).First(&existDevice).Error
	isNewDevice := errors.Is(err, gorm.ErrRecordNotFound)

	if isNewDevice {
		// 新设备，检查设备数量
		var deviceCount int64
		global.DB.Model(&entity.SSODevice{}).
			Where("user_uuid = ? AND app_id = ? AND status = 1", user.UUID, app.ID).
			Count(&deviceCount)

		if int(deviceCount) >= app.MaxDevices {
			return nil, fmt.Errorf("设备数已达上限(%d台)，请先在设备管理中移除其他设备", app.MaxDevices)
		}

		// 创建新设备
		device := entity.SSODevice{
			UserUUID:     user.UUID,
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
		user.UUID,
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
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", user.UUID.String(), deviceID)
	global.Redis.Set(refreshTokenKey, refreshToken, refreshTokenDuration)

	// 记录登录日志
	loginLog := entity.SSOLoginLog{
		UserUUID:  user.UUID,
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
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", claims.UserUUID.String(), claims.DeviceID)
	storedToken, err := global.Redis.Get(refreshTokenKey).Result()
	if err != nil || storedToken != req.RefreshToken {
		return nil, errors.New("refresh_token已被撤销，请重新登录")
	}

	// 检查用户状态
	var user entity.SSOUser
	if err := global.DB.Where("uuid = ?", claims.UserUUID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用或注销")
	}

	// 检查设备状态
	var device entity.SSODevice
	err = global.DB.Where("user_uuid = ? AND device_id = ?", user.UUID, claims.DeviceID).First(&device).Error
	if err != nil || device.Status != 1 {
		return nil, errors.New("设备已被移除")
	}

	// 生成新的AccessToken
	accessTokenDuration, _ := utils.ParseDuration(global.Config.JWT.AccessTokenExpiryTime)
	accessToken, err := jwt.CreateAccessToken(
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
		user.UUID,
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
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", claims.UserUUID.String(), claims.DeviceID)
	global.Redis.Del(refreshTokenKey)

	// 更新设备状态（基于 UUID）
	global.DB.Model(&entity.SSODevice{}).
		Where("user_uuid = ? AND device_id = ?", claims.UserUUID, claims.DeviceID).
		Update("status", 0)

	// 查询应用ID（用于日志）
	var logAppID uint
	if app, err := s.getAppByKey(claims.AppID); err == nil {
		logAppID = app.ID
	}

	// 记录登出日志
	loginLog := entity.SSOLoginLog{
		UserUUID: claims.UserUUID,
		AppID:    logAppID,
		Action:   "logout",
		DeviceID: claims.DeviceID,
		Status:   1,
		Message:  "登出成功",
	}
	global.DB.Create(&loginLog)

	return nil
}

// GetUserInfoByUUID 根据UUID获取用户信息
func (s *AuthService) GetUserInfoByUUID(userUUID uuid.UUID) (*response.UserInfo, error) {
	var user entity.SSOUser
	if err := global.DB.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &response.UserInfo{
		UUID:           user.UUID.String(),
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Email:          user.Email,
		Address:        user.Address,
		Signature:      user.Signature,
		RegisterSource: int(user.RegisterSource),
	}, nil
}

// UpdateUserInfoByUUID 更新用户信息
func (s *AuthService) UpdateUserInfoByUUID(userUUID uuid.UUID, req request.UpdateUserInfoRequest) error {
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

	return global.DB.Model(&entity.SSOUser{}).Where("uuid = ?", userUUID).Updates(updates).Error
}

// UpdatePasswordByUUID 修改密码
func (s *AuthService) UpdatePasswordByUUID(userUUID uuid.UUID, req request.UpdatePasswordRequest) error {
	var user entity.SSOUser
	if err := global.DB.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
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

// GetUserByUUID 根据UUID获取用户信息（服务间调用）
func (s *AuthService) GetUserByUUID(userUUIDStr string) (*response.UserInfo, error) {
	// 解析UUID
	userUUID, err := uuid.FromString(userUUIDStr)
	if err != nil {
		return nil, errors.New("无效的UUID格式")
	}

	// 查询用户
	var user entity.SSOUser
	if err := global.DB.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &response.UserInfo{
		UUID:           user.UUID.String(),
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Email:          user.Email,
		Address:        user.Address,
		Signature:      user.Signature,
		RegisterSource: int(user.RegisterSource), // 注册来源
	}, nil
}

// QQLogin QQ登录
func (s *AuthService) QQLogin(req request.QQLoginRequest, ipAddress, userAgent string) (*response.TokenResponse, error) {
	// 获取QQ access token
	qqService := &QQService{}
	accessTokenResp, err := qqService.GetAccessTokenByCode(req.Code)
	if err != nil || accessTokenResp.Openid == "" {
		return nil, errors.New("获取QQ授权失败")
	}

	// 获取QQ用户信息
	qqUserInfoResp, err := qqService.GetUserInfoByAccessTokenAndOpenid(accessTokenResp.AccessToken, accessTokenResp.Openid)
	if err != nil {
		return nil, errors.New("获取QQ用户信息失败")
	}

	// 将QQ用户信息转换为map
	qqUserInfo := map[string]interface{}{
		"nickname":       qqUserInfoResp.Nickname,
		"figureurl_qq_2": qqUserInfoResp.FigureurlQQ2,
	}

	// 调用QQ登录服务
	deviceID := req.DeviceID
	if deviceID == "" {
		deviceID = uuid.Must(uuid.NewV4()).String()
	}

	return qqService.QQLogin(
		accessTokenResp.Openid,
		req.AppID,
		deviceID,
		req.DeviceName,
		req.DeviceType,
		ipAddress,
		userAgent,
		qqUserInfo,
	)
}

// SendEmailVerificationCode 发送邮箱验证码
func (s *AuthService) SendEmailVerificationCode(email string) error {
	// 生成6位验证码
	verificationCode := utils.GenerateVerificationCode(6)
	expireTime := 5 * time.Minute

	// 存储到Redis，key格式：email_verification_code:{email}
	key := fmt.Sprintf("email_verification_code:%s", email)
	if err := global.Redis.Set(key, verificationCode, expireTime).Err(); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	// 发送邮件
	subject := "您的邮箱验证码"
	body := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #667eea;">邮箱验证码</h2>
		<p>亲爱的用户，</p>
		<p>您正在使用邮箱验证码功能，验证码如下：</p>
		<div style="background: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0; border-radius: 8px;">
			<strong style="font-size: 24px; color: #667eea; letter-spacing: 4px;">%s</strong>
		</div>
		<p>该验证码在 <strong>5 分钟</strong>内有效，请尽快使用。</p>
		<p>如果您没有请求此验证码，请忽略此邮件。</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
	</div>
</body>
</html>`, verificationCode)

	if err := utils.SendEmail(email, subject, body); err != nil {
		// 即使邮件发送失败，验证码已经存储到Redis，记录错误但不返回错误
		// 这样用户仍然可以使用验证码（如果邮件发送成功的话）
		global.Log.Warn("发送邮件失败", zap.Error(err))
		// 可以选择是否返回错误，这里选择不返回，因为验证码已经存储
		// return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// ForgotPassword 忘记密码
func (s *AuthService) ForgotPassword(req request.ForgotPasswordRequest) error {
	// 从Redis获取验证码
	key := fmt.Sprintf("email_verification_code:%s", req.Email)
	storedCode, err := global.Redis.Get(key).Result()
	if err != nil {
		return errors.New("验证码已过期或不存在")
	}

	// 验证验证码
	if storedCode != req.VerificationCode {
		return errors.New("验证码错误")
	}

	// 查询用户
	var user entity.SSOUser
	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("邮箱不存在")
		}
		return err
	}

	// 哈希新密码
	newPasswordHash, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	if err := global.DB.Model(&user).Update("password_hash", newPasswordHash).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 删除验证码
	global.Redis.Del(key)

	return nil
}
