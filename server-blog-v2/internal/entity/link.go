package entity

import "time"

// Link 友链实体。
type Link struct {
	ID          int64
	Name        string
	URL         string
	Logo        *string
	Description *string
	Sort        int
	IsVisible   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
