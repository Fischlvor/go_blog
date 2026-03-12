package output

// FriendLink 友链。
type FriendLink struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
}
