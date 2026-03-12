package input

// CreateSession 创建会话参数。
type CreateSession struct {
	Title string
}

// ListSessions 会话列表参数。
type ListSessions struct {
	PageParams
}

// SendMessage 发送消息参数。
type SendMessage struct {
	Content string
}

// ListAllSessions 会话列表参数（管理端）。
type ListAllSessions struct {
	PageParams
	Keyword  *KeywordParams
	UserUUID *string
}

// ListAllMessages 消息列表参数（管理端）。
type ListAllMessages struct {
	PageParams
	SessionID *int64
	Role      *string
}
