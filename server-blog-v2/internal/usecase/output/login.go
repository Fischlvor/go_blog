package output

import "time"

// LoginInfo 登录记录信息。
type LoginInfo struct {
	ID          int64     `json:"id"`
	UserUUID    string    `json:"user_uuid"`
	LoginMethod string    `json:"login_method"`
	IP          string    `json:"ip"`
	Address     string    `json:"address"`
	OS          string    `json:"os"`
	DeviceInfo  string    `json:"device_info"`
	BrowserInfo string    `json:"browser_info"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
