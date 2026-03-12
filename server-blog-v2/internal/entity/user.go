package entity

import "time"

const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)

// User 用户实体（从 SSO 同步）。
type User struct {
	ID        int64
	UUID      string
	Nickname  string
	Avatar    string
	Email     *string
	RoleID    int
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
