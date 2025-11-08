package entity

import (
	"auth-service/pkg/global"

	"github.com/gofrs/uuid"
)

// UserAppRelation 用户-应用权限关系表
type UserAppRelation struct {
	global.MODEL
	UserUUID uuid.UUID `json:"user_uuid" gorm:"type:char(36);uniqueIndex:idx_user_app;comment:关联sso_users.uuid"`
	AppID    uint      `json:"app_id" gorm:"uniqueIndex:idx_user_app;comment:关联sso_applications.id"`
	Status   int       `json:"status" gorm:"default:1;comment:1可访问 2禁止访问"`
}

func (UserAppRelation) TableName() string {
	return "user_app_relations"
}
