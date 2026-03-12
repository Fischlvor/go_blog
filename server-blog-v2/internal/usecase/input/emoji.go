package input

// ListEmojis 表情列表参数。
type ListEmojis struct {
	PageParams
	Keyword     string
	GroupKey    string
	SpriteGroup *int
}
