package output

import "time"

// UserProfile 用户资料。
type UserProfile struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
	RoleID    int    `json:"role_id"`
}

// UserAdmin 用户（管理端）。
type UserAdmin struct {
	ID        int64      `json:"id"`
	UUID      string     `json:"uuid"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Avatar    string     `json:"avatar"`
	RoleID    int        `json:"role_id"`
	Freeze    bool       `json:"freeze"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

// UserChart 用户图表数据。
type UserChart struct {
	DateList     []string `json:"date_list"`
	LoginData    []int    `json:"login_data"`
	RegisterData []int    `json:"register_data"`
}
