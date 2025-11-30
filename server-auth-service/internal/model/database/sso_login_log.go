package database

import (
	"time"

	"github.com/gofrs/uuid"
)

// SSOLoginLog 登录日志表
type SSOLoginLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:char(36);index;comment:关联sso_users.uuid"`
	AppID     uint      `json:"app_id" gorm:"index;comment:关联sso_applications.id"`
	Action    string    `json:"action" gorm:"size:20;comment:login/logout/kick"`
	DeviceID  string    `json:"device_id" gorm:"size:100"`
	IPAddress string    `json:"ip_address" gorm:"size:50"`
	UserAgent string    `json:"user_agent" gorm:"size:500"`
	Status    int       `json:"status" gorm:"comment:1成功 0失败"`
	Message   string    `json:"message" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
}

func (SSOLoginLog) TableName() string {
	return "sso_login_logs"
}
