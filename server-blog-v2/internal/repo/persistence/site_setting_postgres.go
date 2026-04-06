package persistence

import (
	"context"

	"gorm.io/gorm"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
)

type siteSettingRow struct {
	ID           int64   `gorm:"column:id"`
	SettingKey   string  `gorm:"column:setting_key"`
	SettingValue string  `gorm:"column:setting_value"`
	SettingType  string  `gorm:"column:setting_type"`
	Description  *string `gorm:"column:description"`
	IsPublic     bool    `gorm:"column:is_public"`
}

type siteSettingRepo struct {
	db *gorm.DB
}

// NewSiteSettingRepo 创建站点配置仓库。
func NewSiteSettingRepo(db *gorm.DB) repo.SiteSettingRepo {
	return &siteSettingRepo{db: db}
}

func (r *siteSettingRepo) ListAll(ctx context.Context) ([]*entity.SiteSetting, error) {
	var rows []siteSettingRow
	err := r.db.WithContext(ctx).
		Table("site_settings").
		Order("id ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return toEntitySiteSettings(rows), nil
}

func (r *siteSettingRepo) ListPublic(ctx context.Context) ([]*entity.SiteSetting, error) {
	var rows []siteSettingRow
	err := r.db.WithContext(ctx).
		Table("site_settings").
		Where("is_public = ?", true).
		Order("id ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return toEntitySiteSettings(rows), nil
}

func (r *siteSettingRepo) GetByKey(ctx context.Context, key string) (*entity.SiteSetting, error) {
	var row siteSettingRow
	err := r.db.WithContext(ctx).
		Table("site_settings").
		Where("setting_key = ?", key).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	items := toEntitySiteSettings([]siteSettingRow{row})
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return items[0], nil
}

func (r *siteSettingRepo) UpdateBatch(ctx context.Context, settings []entity.SiteSetting) error {
	if len(settings) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range settings {
			updates := map[string]interface{}{
				"setting_value": s.SettingValue,
				"setting_type":  s.SettingType,
				"description":   s.Description,
				"is_public":     s.IsPublic,
				"updated_at":    gorm.Expr("NOW()"),
			}

			result := tx.Table("site_settings").Where("setting_key = ?", s.SettingKey).Updates(updates)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		return nil
	})
}

func toEntitySiteSettings(rows []siteSettingRow) []*entity.SiteSetting {
	result := make([]*entity.SiteSetting, len(rows))
	for i, row := range rows {
		result[i] = &entity.SiteSetting{
			ID:           row.ID,
			SettingKey:   row.SettingKey,
			SettingValue: row.SettingValue,
			SettingType:  row.SettingType,
			Description:  row.Description,
			IsPublic:     row.IsPublic,
		}
	}
	return result
}
