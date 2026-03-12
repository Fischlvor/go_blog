package entity

import "time"

const (
	FeedbackStatusPending   = "pending"
	FeedbackStatusProcessed = "processed"
)

// Feedback 反馈实体。
type Feedback struct {
	ID        int64
	UserUUID  *string // 关联用户
	Type      string  // bug, suggestion, other
	Content   string
	Contact   *string
	Status    string
	Reply     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
