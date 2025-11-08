package entity

import (
	"auth-service/pkg/global"
)

// SSOApplication 应用注册表
type SSOApplication struct {
	global.MODEL
	AppKey         string `json:"app_key" gorm:"uniqueIndex;size:50;not null;comment:应用标识如blog/admin"`
	AppName        string `json:"app_name" gorm:"size:100;not null"`
	AppSecret      string `json:"-" gorm:"size:255;not null"`
	RedirectURIs   string `json:"redirect_uris" gorm:"type:text;comment:允许的回调地址，逗号分隔"`
	MaxDevices     int    `json:"max_devices" gorm:"default:5;comment:单应用最大设备数"`
	AllowedOrigins string `json:"allowed_origins" gorm:"type:text;comment:CORS白名单，逗号分隔"`
	Status         int    `json:"status" gorm:"default:1;comment:1启用 0禁用"`
}

func (SSOApplication) TableName() string {
	return "sso_applications"
}
