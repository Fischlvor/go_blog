package entity

import "time"

// Advertisement 广告实体。
type Advertisement struct {
	ID        int64
	AdName    string
	AdImage   string
	AdLink    string
	AdType    int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
