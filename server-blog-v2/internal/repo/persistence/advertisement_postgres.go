package persistence

import (
	"context"

	"gorm.io/gorm"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
)

type advertisementRepo struct {
	db *gorm.DB
}

// NewAdvertisementRepo 创建 Advertisement 仓库。
func NewAdvertisementRepo(db *gorm.DB) repo.AdvertisementRepo {
	return &advertisementRepo{db: db}
}

func (r *advertisementRepo) ListActive(ctx context.Context) ([]*entity.Advertisement, int64, error) {
	var models []struct {
		ID       int64  `gorm:"column:id"`
		AdName   string `gorm:"column:ad_name"`
		AdImage  string `gorm:"column:ad_image"`
		AdLink   string `gorm:"column:ad_link"`
		AdType   int    `gorm:"column:ad_type"`
		IsActive bool   `gorm:"column:is_active"`
	}

	var total int64
	db := r.db.WithContext(ctx).Table("advertisements").Where("is_active = ?", true)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entity.Advertisement, len(models))
	for i, m := range models {
		result[i] = &entity.Advertisement{
			ID:       m.ID,
			AdName:   m.AdName,
			AdImage:  m.AdImage,
			AdLink:   m.AdLink,
			AdType:   m.AdType,
			IsActive: m.IsActive,
		}
	}

	return result, total, nil
}

func (r *advertisementRepo) List(ctx context.Context, offset, limit int, title *string) ([]*entity.Advertisement, int64, error) {
	var models []struct {
		ID       int64  `gorm:"column:id"`
		AdName   string `gorm:"column:ad_name"`
		AdImage  string `gorm:"column:ad_image"`
		AdLink   string `gorm:"column:ad_link"`
		AdType   int    `gorm:"column:ad_type"`
		IsActive bool   `gorm:"column:is_active"`
	}

	db := r.db.WithContext(ctx).Table("advertisements")

	if title != nil && *title != "" {
		db = db.Where("ad_name LIKE ?", "%"+*title+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entity.Advertisement, len(models))
	for i, m := range models {
		result[i] = &entity.Advertisement{
			ID:       m.ID,
			AdName:   m.AdName,
			AdImage:  m.AdImage,
			AdLink:   m.AdLink,
			AdType:   m.AdType,
			IsActive: m.IsActive,
		}
	}

	return result, total, nil
}
