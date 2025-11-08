package entity

import (
	"auth-service/pkg/global"
	"time"
)

// SSODevice 设备管理表
type SSODevice struct {
	global.MODEL
	UserID       uint      `json:"user_id" gorm:"uniqueIndex:idx_user_app_device;comment:关联sso_users.id"`
	AppID        uint      `json:"app_id" gorm:"uniqueIndex:idx_user_app_device;comment:关联sso_applications.id"`
	DeviceID     string    `json:"device_id" gorm:"uniqueIndex:idx_user_app_device;size:100;comment:设备标识"`
	DeviceName   string    `json:"device_name" gorm:"size:100;comment:设备名称"`
	DeviceType   string    `json:"device_type" gorm:"size:20;comment:web/ios/android"`
	IPAddress    string    `json:"ip_address" gorm:"size:50"`
	UserAgent    string    `json:"user_agent" gorm:"size:500"`
	LastActiveAt time.Time `json:"last_active_at" gorm:"comment:最后活跃时间"`
	Status       int       `json:"status" gorm:"default:1;index;comment:1在线 0离线"`
}

func (SSODevice) TableName() string {
	return "sso_devices"
}
