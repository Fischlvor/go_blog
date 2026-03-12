package input

// UpdateProfile 更新用户资料参数。
type UpdateProfile struct {
	Nickname string
	Avatar   string
}

// ListUsers 用户列表参数（管理端）。
type ListUsers struct {
	PageParams
	Keyword *KeywordParams
	Status  *string // active, frozen
}
