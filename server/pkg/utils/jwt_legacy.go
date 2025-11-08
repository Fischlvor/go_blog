package utils

import (
	"github.com/gofrs/uuid"
)

// LegacyClaims 兼容旧系统的Claims结构
type LegacyClaims struct {
	UserID   uint      `json:"user_id"`
	UserUUID uuid.UUID `json:"uuid"`
	RoleID   int       `json:"role_id"`
}
