package entity

import (
	"auth-service/pkg/global"
)

// SSOOAuthBinding 第三方登录绑定表
type SSOOAuthBinding struct {
	global.MODEL
	UserID   uint   `json:"user_id" gorm:"index;comment:关联sso_users"`
	Provider string `json:"provider" gorm:"size:20;uniqueIndex:idx_provider_openid;comment:qq/wechat/github"`
	OpenID   string `json:"open_id" gorm:"size:100;uniqueIndex:idx_provider_openid"`
	UnionID  string `json:"union_id" gorm:"size:100"`
}

func (SSOOAuthBinding) TableName() string {
	return "sso_oauth_bindings"
}
