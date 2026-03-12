package request

// ReplyFeedback 回复反馈请求。
type ReplyFeedback struct {
	ID    int64  `json:"id" validate:"required"`
	Reply string `json:"reply" validate:"required,max=1000"`
}
