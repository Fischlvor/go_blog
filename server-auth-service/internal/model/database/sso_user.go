package database

import (
	"auth-service/internal/model/appTypes"
	"auth-service/pkg/global"

	"github.com/gofrs/uuid"
)

// SSOUser SSO用户表
type SSOUser struct {
	global.MODEL
	UUID           uuid.UUID               `json:"uuid" gorm:"type:char(36);uniqueIndex"`
	Username       string                  `json:"username" gorm:"size:50"`
	PasswordHash   string                  `json:"-" gorm:"size:255;column:password_hash"`
	Email          *string                 `json:"email" gorm:"size:100;uniqueIndex;default:null"`
	Phone          string                  `json:"phone" gorm:"size:20;uniqueIndex;default:null"`
	Nickname       string                  `json:"nickname" gorm:"size:50"`
	Avatar         string                  `json:"avatar" gorm:"size:255"`
	Address        string                  `json:"address" gorm:"size:255"`
	Signature      string                  `json:"signature" gorm:"type:text"`
	Status         int                     `json:"status" gorm:"default:1;index;comment:1正常 2禁用 3注销"`
	RegisterSource appTypes.RegisterSource `json:"register_source" gorm:"type:int;default:0;comment:0email 1qq 2wechat 3github"`
	IsSuperAdmin   bool                    `json:"is_super_admin" gorm:"default:false;comment:是否超级管理员"`
}

func (SSOUser) TableName() string {
	return "sso_users"
}
