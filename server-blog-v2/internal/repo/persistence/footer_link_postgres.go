package persistence

import (
	"context"

	"gorm.io/gorm"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
)

type footerLinkRepo struct {
	db *gorm.DB
}

// NewFooterLinkRepo 创建 FooterLink 仓库。
func NewFooterLinkRepo(db *gorm.DB) repo.FooterLinkRepo {
	return &footerLinkRepo{db: db}
}

func (r *footerLinkRepo) List(ctx context.Context) ([]*entity.FooterLink, error) {
	var models []struct {
		Title string `gorm:"column:title"`
		Link  string `gorm:"column:link"`
	}

	err := r.db.WithContext(ctx).
		Table("footer_links").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.FooterLink, len(models))
	for i, m := range models {
		result[i] = &entity.FooterLink{
			Title: m.Title,
			Link:  m.Link,
		}
	}

	return result, nil
}
