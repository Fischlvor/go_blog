package setting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
)

var (
	// ErrRepo Repo 错误哨兵。
	ErrRepo = errors.New("repo")
)

type useCase struct {
	settings repo.SiteSettingRepo
}

// New 创建 Setting UseCase。
func New(settings repo.SiteSettingRepo) usecase.Setting {
	return &useCase{settings: settings}
}

func (u *useCase) GetAllSiteSettings(ctx context.Context) (*output.AllResult[output.SiteSettingDetail], error) {
	settings, err := u.settings.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	items := make([]output.SiteSettingDetail, len(settings))
	for i, s := range settings {
		items[i] = toSiteSettingDetail(s)
	}
	return &output.AllResult[output.SiteSettingDetail]{
		Items: items,
		Total: int64(len(items)),
	}, nil
}

func (u *useCase) GetPublicSiteSettings(ctx context.Context) (*output.AllResult[output.SiteSettingDetail], error) {
	settings, err := u.settings.ListPublic(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	items := make([]output.SiteSettingDetail, len(settings))
	for i, s := range settings {
		items[i] = toSiteSettingDetail(s)
	}
	return &output.AllResult[output.SiteSettingDetail]{
		Items: items,
		Total: int64(len(items)),
	}, nil
}

func (u *useCase) GetSiteSettingByKey(ctx context.Context, key string) (*output.SiteSettingDetail, error) {
	s, err := u.settings.GetByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	detail := toSiteSettingDetail(s)
	return &detail, nil
}

func (u *useCase) GetWebsiteSettingsMap(ctx context.Context) (map[string]string, error) {
	settings, err := u.settings.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.SettingKey] = s.SettingValue
	}
	return result, nil
}

func (u *useCase) UpdateSiteSettings(ctx context.Context, settings []input.UpsertSiteSetting) error {
	items := make([]entity.SiteSetting, len(settings))
	for i, s := range settings {
		var desc *string
		if s.Description != "" {
			d := s.Description
			desc = &d
		}
		settingType := s.SettingType
		if settingType == "" {
			settingType = "string"
		}
		items[i] = entity.SiteSetting{
			SettingKey:   s.SettingKey,
			SettingValue: s.SettingValue,
			SettingType:  settingType,
			Description:  desc,
			IsPublic:     s.IsPublic,
		}
	}
	if err := u.settings.UpdateBatch(ctx, items); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

func toSiteSettingDetail(s *entity.SiteSetting) output.SiteSettingDetail {
	description := ""
	if s.Description != nil {
		description = *s.Description
	}
	return output.SiteSettingDetail{
		ID:           s.ID,
		SettingKey:   s.SettingKey,
		SettingValue: s.SettingValue,
		SettingType:  s.SettingType,
		Description:  description,
		IsPublic:     s.IsPublic,
		CreatedAt:    formatTime(s.CreatedAt),
		UpdatedAt:    formatTime(s.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
