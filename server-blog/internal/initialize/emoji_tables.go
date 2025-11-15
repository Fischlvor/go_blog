package initialize

import (
	"server/internal/model"
	"server/pkg/global"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

// InitEmojiTables 初始化Emoji相关表
func InitEmojiTables() {
	// 创建Emoji相关表
	err := global.DB.AutoMigrate(
		&model.Emoji{},
		&model.EmojiGroup{},
		&model.EmojiSprite{},
		&model.EmojiTask{},
	)
	if err != nil {
		global.Log.Error("创建Emoji表失败", zap.Error(err))
		return
	}
	global.Log.Info("✓ Emoji表创建成功")
}

// InitDefaultEmojiGroups 初始化默认emoji组
func InitDefaultEmojiGroups() {
	// 检查是否已存在默认组
	var count int64
	global.DB.Model(&model.EmojiGroup{}).Where("group_name = ?", "系统基础表情").Count(&count)
	if count > 0 {
		global.Log.Info("✓ 默认emoji组已存在，跳过初始化")
		return
	}

	// 管理员UUID
	adminUUID, err := uuid.FromString("37395c61-a2ec-464e-9567-ce6fa92630f7")
	if err != nil {
		global.Log.Error("解析管理员UUID失败", zap.Error(err))
		return
	}

	// 创建默认组
	defaultGroup := &model.EmojiGroup{
		GroupName:   "系统基础表情",
		GroupKey:    "system_1_base",
		Description: "系统内置的基础表情包",
		SortOrder:   1,
		EmojiCount:  0,
		Status:      1,
		CreatedBy:   adminUUID,
	}

	if err := global.DB.Create(defaultGroup).Error; err != nil {
		global.Log.Error("创建默认emoji组失败", zap.Error(err))
		return
	}

	global.Log.Info("✓ 默认emoji组创建成功")
}
