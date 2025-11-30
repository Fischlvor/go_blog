package database

import (
	"auth-service/pkg/global"

	"github.com/gofrs/uuid"
)

// SSOOAuthBinding 第三方登录绑定表
type SSOOAuthBinding struct {
	global.MODEL
	UserUUID uuid.UUID `json:"user_uuid" gorm:"type:char(36);index;comment:关联sso_users.uuid"`
	Provider string    `json:"provider" gorm:"size:20;uniqueIndex:idx_provider_openid;comment:qq/wechat/github"`
	OpenID   string    `json:"open_id" gorm:"size:100;uniqueIndex:idx_provider_openid"`
	UnionID  string    `json:"union_id" gorm:"size:100"`
}

func (SSOOAuthBinding) TableName() string {
	return "sso_oauth_bindings"
}
