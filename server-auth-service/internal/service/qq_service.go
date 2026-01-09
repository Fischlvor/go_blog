package service

import (
	types "auth-service/internal/model/appTypes"
	database "auth-service/internal/model/database"
	"auth-service/internal/model/other"
	"auth-service/internal/model/response"
	"auth-service/pkg/global"
	"auth-service/pkg/jwt"
	"auth-service/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type QQService struct{}

// GetAccessTokenByCode 通过Authorization Code获取Access Token
func (s *QQService) GetAccessTokenByCode(code string) (other.AccessTokenResponse, error) {
	data := other.AccessTokenResponse{}
	clientID := global.Config.QQ.AppID
	clientSecret := global.Config.QQ.AppKey
	redirectUri := global.Config.QQ.RedirectURI
	urlStr := "https://graph.qq.com/oauth2.0/token"
	method := "GET"
	params := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectUri,
		"fmt":           "json",
		"need_openid":   "1",
	}
	res, err := utils.HttpRequest(urlStr, method, nil, params, nil)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return data, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// GetUserInfoByAccessTokenAndOpenid 获取登录用户信息
func (s *QQService) GetUserInfoByAccessTokenAndOpenid(accessToken, openID string) (other.UserInfoResponse, error) {
	data := other.UserInfoResponse{}
	oauthConsumerKey := global.Config.QQ.AppID
	urlStr := "https://graph.qq.com/user/get_user_info"
	method := "GET"
	params := map[string]string{
		"access_token":       accessToken,
		"oauth_consumer_key": oauthConsumerKey,
		"openid":             openID,
	}
	res, err := utils.HttpRequest(urlStr, method, nil, params, nil)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return data, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// getAppByKey 通过app_key获取应用信息
func (s *QQService) getAppByKey(appKey string) (*database.SSOApplication, error) {
	var app database.SSOApplication
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
	var oauthBinding database.SSOOAuthBinding
	err := global.DB.Where("provider = ? AND open_id = ?", "qq", openID).First(&oauthBinding).Error

	var user database.SSOUser
	var isNewUser bool

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 新用户，创建账号
		isNewUser = true

		// 从QQ获取昵称和头像
		nickname := "QQ用户"
		avatar := "https://image.hsk423.cn/blog/aaca0f5eb4d2d98a6ce6dffa99f8254b-20251108151238.jpg"
		if name, ok := qqUserInfo["nickname"].(string); ok && name != "" {
			nickname = name
		}
		if figureurl, ok := qqUserInfo["figureurl_qq_2"].(string); ok && figureurl != "" {
			avatar = figureurl
		}

		user = database.SSOUser{
			UUID:           uuid.Must(uuid.NewV4()),
			Username:       nickname,
			PasswordHash:   nil, // QQ登录无密码
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
		oauthBinding = database.SSOOAuthBinding{
			UserUUID: user.UUID,
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
		userAppRelation := database.UserAppRelation{
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
		// 老用户，查询用户信息（通过 UUID）
		if err := global.DB.Where("uuid = ?", oauthBinding.UserUUID).First(&user).Error; err != nil {
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
	var userAppRelation database.UserAppRelation
	err = global.DB.Where("user_uuid = ? AND app_id = ?", user.UUID, app.ID).First(&userAppRelation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 自动授权
			userAppRelation = database.UserAppRelation{
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
	var existDevice database.SSODevice
	err = global.DB.Where("user_uuid = ? AND app_id = ? AND device_id = ?", user.UUID, app.ID, deviceID).First(&existDevice).Error
	isNewDevice := errors.Is(err, gorm.ErrRecordNotFound)

	if isNewDevice {
		// 新设备，检查设备数量限制并自动踢出最早的设备
		authService := &AuthService{}
		err = authService.handleDeviceLimit(user.UUID, app.ID, app.MaxDevices)
		if err != nil {
			return nil, fmt.Errorf("处理设备限制失败: %w", err)
		}
	}

	device := database.SSODevice{
		UserUUID:     user.UUID,
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

		// 记录新设备QQ登录日志
		authService := &AuthService{}
		authService.LogAction(user.UUID, app.ID, "qq_login", deviceID, "QQ新设备登录成功", 1)
	} else {
		global.DB.Model(&existDevice).Updates(map[string]interface{}{
			"last_active_at": time.Now(),
			"ip_address":     ipAddress,
			"user_agent":     userAgent,
			"status":         1,
		})

		// 记录现有设备QQ登录日志
		authService := &AuthService{}
		authService.LogAction(user.UUID, app.ID, "qq_login", deviceID, "QQ设备登录成功", 1)
	}

	// 生成Token
	accessTokenDuration, _ := utils.ParseDuration(global.Config.JWT.AccessTokenExpiryTime)
	refreshTokenDuration, _ := utils.ParseDuration(global.Config.JWT.RefreshTokenExpiryTime)

	accessToken, _ := jwt.CreateAccessToken(
		user.UUID,
		appID,
		deviceID,
		accessTokenDuration,
		global.Config.JWT.Issuer,
		global.RSAPrivateKey,
	)

	refreshToken, _ := jwt.CreateRefreshToken(
		user.UUID,
		appID,
		deviceID,
		refreshTokenDuration,
		global.Config.JWT.Issuer,
		global.RSAPrivateKey,
	)

	// 存储RefreshToken
	refreshTokenKey := fmt.Sprintf("refresh_token:%s:%s", user.UUID.String(), deviceID)
	global.Redis.Set(refreshTokenKey, refreshToken, refreshTokenDuration)

	// 记录登录日志
	loginLog := database.SSOLoginLog{
		UserUUID:  user.UUID,
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
			UUID:      user.UUID.String(),
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Email:     user.Email,
			Address:   user.Address,
			Signature: user.Signature,
		},
	}, nil
}
