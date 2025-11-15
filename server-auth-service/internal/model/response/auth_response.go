package response

import (
	"time"
)

// TokenResponse Token响应
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	UserInfo     *UserInfo `json:"user_info,omitempty"`
}

// UserInfo 用户信息
type UserInfo struct {
	UUID           string `json:"uuid"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	Signature      string `json:"signature"`
	RegisterSource int    `json:"register_source"` // 0:email 1:qq 2:wechat 3:github
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceID     string    `json:"device_id"`
	DeviceName   string    `json:"device_name"`
	DeviceType   string    `json:"device_type"`
	IPAddress    string    `json:"ip_address"`
	LastActiveAt time.Time `json:"last_active_at"`
	Status       int       `json:"status"`
	IsCurrent    bool      `json:"is_current"`
}

// LoginLogInfo 登录日志信息
type LoginLogInfo struct {
	ID        uint      `json:"id"`
	AppID     string    `json:"app_id"`
	Action    string    `json:"action"`
	DeviceID  string    `json:"device_id"`
	IPAddress string    `json:"ip_address"`
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// ApplicationInfo 应用信息
type ApplicationInfo struct {
	AppID      string    `json:"app_id"`
	AppName    string    `json:"app_name"`
	MaxDevices int       `json:"max_devices"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// UserAppPermission 用户应用权限
type UserAppPermission struct {
	AppID   string `json:"app_id"`
	AppName string `json:"app_name"`
	Status  int    `json:"status"`
}
