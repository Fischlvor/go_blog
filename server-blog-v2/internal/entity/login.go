package entity

import "time"

// Login 登录记录实体。
type Login struct {
	ID          int64
	UserUUID    string
	LoginMethod string
	IP          string
	Address     string
	OS          string
	DeviceInfo  string
	BrowserInfo string
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
