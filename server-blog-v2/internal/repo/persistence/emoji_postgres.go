package persistence

import (
	"context"

	"gorm.io/gorm"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
)

type emojiRepo struct {
	db *gorm.DB
}

// NewEmojiRepo 创建 Emoji 仓库。
func NewEmojiRepo(db *gorm.DB) repo.EmojiRepo {
	return &emojiRepo{db: db}
}

func (r *emojiRepo) ListActive(ctx context.Context) ([]*entity.Emoji, error) {
	var models []struct {
		ID              int64  `gorm:"column:id"`
		Key             string `gorm:"column:key"`
		Filename        string `gorm:"column:filename"`
		GroupKey        string `gorm:"column:group_key"`
		SpriteGroup     int    `gorm:"column:sprite_group"`
		SpritePositionX int    `gorm:"column:sprite_position_x"`
		SpritePositionY int    `gorm:"column:sprite_position_y"`
		FileSize        int    `gorm:"column:file_size"`
		CdnURL          string `gorm:"column:cdn_url"`
		Status          int8   `gorm:"column:status"`
	}

	err := r.db.WithContext(ctx).
		Table("emojis").
		Where("status = ?", entity.EmojiStatusActive).
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Emoji, len(models))
	for i, m := range models {
		result[i] = &entity.Emoji{
			ID:              m.ID,
			Key:             m.Key,
			Filename:        m.Filename,
			GroupKey:        m.GroupKey,
			SpriteGroup:     m.SpriteGroup,
			SpritePositionX: m.SpritePositionX,
			SpritePositionY: m.SpritePositionY,
			FileSize:        m.FileSize,
			CdnURL:          m.CdnURL,
			Status:          m.Status,
		}
	}

	return result, nil
}

func (r *emojiRepo) ListGroups(ctx context.Context) ([]*entity.EmojiGroup, error) {
	var models []struct {
		ID            int64  `gorm:"column:id"`
		GroupName     string `gorm:"column:group_name"`
		GroupKey      string `gorm:"column:group_key"`
		Description   string `gorm:"column:description"`
		SortOrder     int    `gorm:"column:sort_order"`
		EmojiCount    int    `gorm:"column:emoji_count"`
		Status        int8   `gorm:"column:status"`
		SpriteConfURL string `gorm:"column:sprite_conf_url"`
		CreatedAt     string `gorm:"column:created_at"`
		UpdatedAt     string `gorm:"column:updated_at"`
	}

	err := r.db.WithContext(ctx).
		Table("emoji_groups").
		Order("sort_order ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.EmojiGroup, len(models))
	for i, m := range models {
		result[i] = &entity.EmojiGroup{
			ID:            m.ID,
			GroupName:     m.GroupName,
			GroupKey:      m.GroupKey,
			Description:   m.Description,
			SortOrder:     m.SortOrder,
			EmojiCount:    m.EmojiCount,
			Status:        m.Status,
			SpriteConfURL: m.SpriteConfURL,
		}
	}

	return result, nil
}

func (r *emojiRepo) ListActiveByGroupKeys(ctx context.Context, groupKeys []string) ([]*entity.Emoji, error) {
	var models []struct {
		ID              int64  `gorm:"column:id"`
		Key             string `gorm:"column:key"`
		Filename        string `gorm:"column:filename"`
		GroupKey        string `gorm:"column:group_key"`
		SpriteGroup     int    `gorm:"column:sprite_group"`
		SpritePositionX int    `gorm:"column:sprite_position_x"`
		SpritePositionY int    `gorm:"column:sprite_position_y"`
		FileSize        int    `gorm:"column:file_size"`
		CdnURL          string `gorm:"column:cdn_url"`
		Status          int8   `gorm:"column:status"`
	}

	db := r.db.WithContext(ctx).
		Table("emojis").
		Where("status = ?", entity.EmojiStatusActive).
		Order("group_key ASC, id ASC")

	if len(groupKeys) > 0 {
		db = db.Where("group_key IN ?", groupKeys)
	}

	err := db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Emoji, len(models))
	for i, m := range models {
		result[i] = &entity.Emoji{
			ID:              m.ID,
			Key:             m.Key,
			Filename:        m.Filename,
			GroupKey:        m.GroupKey,
			SpriteGroup:     m.SpriteGroup,
			SpritePositionX: m.SpritePositionX,
			SpritePositionY: m.SpritePositionY,
			FileSize:        m.FileSize,
			CdnURL:          m.CdnURL,
			Status:          m.Status,
		}
	}

	return result, nil
}

func (r *emojiRepo) GetGroupByKey(ctx context.Context, groupKey string) (*entity.EmojiGroup, error) {
	var m struct {
		ID            int64  `gorm:"column:id"`
		GroupName     string `gorm:"column:group_name"`
		GroupKey      string `gorm:"column:group_key"`
		Description   string `gorm:"column:description"`
		SortOrder     int    `gorm:"column:sort_order"`
		EmojiCount    int    `gorm:"column:emoji_count"`
		Status        int8   `gorm:"column:status"`
		SpriteConfURL string `gorm:"column:sprite_conf_url"`
	}

	err := r.db.WithContext(ctx).
		Table("emoji_groups").
		Where("group_key = ?", groupKey).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return &entity.EmojiGroup{
		ID:            m.ID,
		GroupName:     m.GroupName,
		GroupKey:      m.GroupKey,
		Description:   m.Description,
		SortOrder:     m.SortOrder,
		EmojiCount:    m.EmojiCount,
		Status:        m.Status,
		SpriteConfURL: m.SpriteConfURL,
	}, nil
}

func (r *emojiRepo) UpdateGroupSpriteConfURL(ctx context.Context, groupKey, confURL string) error {
	return r.db.WithContext(ctx).
		Table("emoji_groups").
		Where("group_key = ?", groupKey).
		Update("sprite_conf_url", confURL).Error
}

func (r *emojiRepo) List(ctx context.Context, offset, limit int, keyword, groupKey string, spriteGroup *int) ([]*entity.Emoji, int64, error) {
	var models []struct {
		ID              int64  `gorm:"column:id"`
		Key             string `gorm:"column:key"`
		Filename        string `gorm:"column:filename"`
		GroupKey        string `gorm:"column:group_key"`
		SpriteGroup     int    `gorm:"column:sprite_group"`
		SpritePositionX int    `gorm:"column:sprite_position_x"`
		SpritePositionY int    `gorm:"column:sprite_position_y"`
		FileSize        int    `gorm:"column:file_size"`
		CdnURL          string `gorm:"column:cdn_url"`
		Status          int8   `gorm:"column:status"`
		UploadTime      string `gorm:"column:upload_time"`
		UpdatedAt       string `gorm:"column:updated_at"`
	}

	db := r.db.WithContext(ctx).Table("emojis")

	if keyword != "" {
		db = db.Where("key LIKE ? OR filename LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if groupKey != "" {
		db = db.Where("group_key = ?", groupKey)
	}
	if spriteGroup != nil {
		db = db.Where("sprite_group = ?", *spriteGroup)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.Emoji, len(models))
	for i, m := range models {
		result[i] = &entity.Emoji{
			ID:              m.ID,
			Key:             m.Key,
			Filename:        m.Filename,
			GroupKey:        m.GroupKey,
			SpriteGroup:     m.SpriteGroup,
			SpritePositionX: m.SpritePositionX,
			SpritePositionY: m.SpritePositionY,
			FileSize:        m.FileSize,
			CdnURL:          m.CdnURL,
			Status:          m.Status,
		}
	}

	return result, total, nil
}

type emojiSpriteRepo struct {
	db *gorm.DB
}

// NewEmojiSpriteRepo 创建 EmojiSprite 仓库。
func NewEmojiSpriteRepo(db *gorm.DB) repo.EmojiSpriteRepo {
	return &emojiSpriteRepo{db: db}
}

func (r *emojiSpriteRepo) ListActive(ctx context.Context) ([]*entity.EmojiSprite, error) {
	var models []struct {
		ID          int64  `gorm:"column:id"`
		SpriteGroup int    `gorm:"column:sprite_group"`
		Filename    string `gorm:"column:filename"`
		CdnURL      string `gorm:"column:cdn_url"`
		Width       int    `gorm:"column:width"`
		Height      int    `gorm:"column:height"`
		EmojiCount  int    `gorm:"column:emoji_count"`
		FileSize    int    `gorm:"column:file_size"`
		Status      int8   `gorm:"column:status"`
	}

	err := r.db.WithContext(ctx).
		Table("emoji_sprites").
		Where("status = ?", entity.EmojiSpriteStatusActive).
		Order("sprite_group ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.EmojiSprite, len(models))
	for i, m := range models {
		result[i] = &entity.EmojiSprite{
			ID:          m.ID,
			SpriteGroup: m.SpriteGroup,
			Filename:    m.Filename,
			CdnURL:      m.CdnURL,
			Width:       m.Width,
			Height:      m.Height,
			EmojiCount:  m.EmojiCount,
			FileSize:    m.FileSize,
			Status:      m.Status,
		}
	}

	return result, nil
}

func (r *emojiSpriteRepo) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM emoji_sprites").Error
}

func (r *emojiSpriteRepo) CreateBatch(ctx context.Context, sprites []*entity.EmojiSprite) error {
	if len(sprites) == 0 {
		return nil
	}

	models := make([]map[string]interface{}, len(sprites))
	for i, s := range sprites {
		models[i] = map[string]interface{}{
			"sprite_group": s.SpriteGroup,
			"filename":     s.Filename,
			"cdn_url":      s.CdnURL,
			"width":        s.Width,
			"height":       s.Height,
			"emoji_count":  s.EmojiCount,
			"file_size":    s.FileSize,
			"status":       s.Status,
		}
	}

	return r.db.WithContext(ctx).Table("emoji_sprites").Create(models).Error
}
