package flag

import (
	"server/internal/initialize"
)

// Emoji 初始化Emoji相关的数据库表和默认数据
func Emoji() error {
	// 初始化Emoji相关表
	initialize.InitEmojiTables()

	// 初始化默认emoji组
	initialize.InitDefaultEmojiGroups()

	return nil
}
