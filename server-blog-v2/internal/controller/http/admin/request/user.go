package request

// FreezeUser 冻结用户请求。
type FreezeUser struct {
	UserUUID string `json:"user_uuid" validate:"required"`
}

// UnfreezeUser 解冻用户请求。
type UnfreezeUser struct {
	UserUUID string `json:"user_uuid" validate:"required"`
}
