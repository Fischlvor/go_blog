package output

import "time"

// Feedback 反馈。
type Feedback struct {
	ID        int64     `json:"id"`
	UserUUID  string    `json:"user_uuid,omitempty"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Contact   string    `json:"contact,omitempty"`
	Status    string    `json:"status"`
	Reply     string    `json:"reply,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
