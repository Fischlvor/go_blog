package website

import (
	"context"
	"encoding/json"
	"time"

	"server-blog-v2/config"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/output"
	"server-blog-v2/internal/usecase/urlutil"
	"server-blog-v2/pkg/calendar"
	"server-blog-v2/pkg/hotsearch"
	"server-blog-v2/pkg/redis"
)

type useCase struct {
	cfg         *config.Config
	redis       redis.Client
	footerLinks repo.FooterLinkRepo
	settings    repo.SiteSettingRepo
}

// New 创建 Website UseCase。
func New(cfg *config.Config, redis redis.Client, footerLinks repo.FooterLinkRepo, settings repo.SiteSettingRepo) usecase.Website {
	return &useCase{
		cfg:         cfg,
		redis:       redis,
		footerLinks: footerLinks,
		settings:    settings,
	}
}

func (u *useCase) GetInfo(ctx context.Context) *output.WebsiteInfo {
	w := u.cfg.Website
	values := u.loadWebsiteSettings(ctx)
	return &output.WebsiteInfo{
		Avatar:          settingValueWithKey(values, urlutil.ResolveImageURL(u.cfg, firstNonEmpty(values["profile.avatar"], values["website.avatar"], w.Avatar)), "profile.avatar", "website.avatar"),
		Title:           settingValueWithKey(values, firstNonEmpty(values["website.title"], w.Title), "website.title"),
		Description:     settingValueWithKey(values, firstNonEmpty(values["website.description"], w.Description), "website.description"),
		ProfileIntro:    settingValueWithKey(values, firstNonEmpty(values["profile.intro"]), "profile.intro"),
		TechStack:       settingValueWithKey(values, firstNonEmpty(values["profile.tech_stack"]), "profile.tech_stack"),
		WorkExperiences: settingValueWithKey(values, firstNonEmpty(values["profile.work_experiences"]), "profile.work_experiences"),
		Version:         settingValueWithKey(values, firstNonEmpty(values["website.version"], w.Version), "website.version"),
		CreatedAt:       settingValueWithKey(values, firstNonEmpty(values["website.created_at"], w.CreatedAt), "website.created_at"),
		ICPFiling:       settingValueWithKey(values, firstNonEmpty(values["website.icp_filing"], w.ICPFiling), "website.icp_filing"),
		BilibiliURL:     settingValueWithKey(values, firstNonEmpty(values["profile.bilibili_url"], w.BilibiliURL), "profile.bilibili_url"),
		GithubURL:       settingValueWithKey(values, firstNonEmpty(values["profile.github_url"], w.GithubURL), "profile.github_url"),
		SteamURL:        settingValueWithKey(values, firstNonEmpty(values["profile.steam_url"]), "profile.steam_url"),
		Name:            settingValueWithKey(values, firstNonEmpty(values["profile.name"], w.Name), "profile.name"),
		Job:             settingValueWithKey(values, firstNonEmpty(values["profile.job"], w.Job), "profile.job"),
		Address:         settingValueWithKey(values, firstNonEmpty(values["profile.address"], w.Address), "profile.address"),
		Email:           settingValueWithKey(values, firstNonEmpty(values["profile.email"], w.Email), "profile.email"),
	}
}

func (u *useCase) GetCarousel(ctx context.Context) ([]string, error) {
	// TODO: 从数据库获取轮播图
	// 暂时返回默认值
	return []string{
		"/image/carousel_1.jpg",
		"/image/carousel_2.jpg",
		"/image/carousel_3.jpg",
		"/image/carousel_4.jpg",
	}, nil
}

func (u *useCase) GetNews(ctx context.Context, source string) (*output.HotSearchData, error) {
	// 先从 Redis 缓存获取
	result, err := u.redis.Get(ctx, source)
	if err == nil && result != "" {
		var data output.HotSearchData
		if err := json.Unmarshal([]byte(result), &data); err == nil {
			return &data, nil
		}
	}

	// 缓存未命中，从源获取
	src := hotsearch.NewSource(source)
	if src == nil {
		return &output.HotSearchData{Source: source, HotList: []output.HotItem{}}, nil
	}

	hsData, err := src.GetHotSearchData(30)
	if err != nil {
		return nil, err
	}

	// 转换为 output 类型
	data := &output.HotSearchData{
		Source:     hsData.Source,
		UpdateTime: hsData.UpdateTime,
		HotList:    make([]output.HotItem, len(hsData.HotList)),
	}
	for i, item := range hsData.HotList {
		data.HotList[i] = output.HotItem{
			Index:       item.Index,
			Title:       item.Title,
			Description: item.Description,
			Image:       item.Image,
			Popularity:  item.Popularity,
			URL:         item.URL,
		}
	}

	// 缓存到 Redis，1 小时过期
	bytes, _ := json.Marshal(data)
	_ = u.redis.Set(ctx, source, string(bytes), time.Hour)

	return data, nil
}

func (u *useCase) GetCalendar(ctx context.Context) (*output.CalendarData, error) {
	dateStr := time.Now().Format("2006/0102")
	cacheKey := "calendar-" + dateStr

	// 先从 Redis 缓存获取
	result, err := u.redis.Get(ctx, cacheKey)
	if err == nil && result != "" {
		var cal output.CalendarData
		if err := json.Unmarshal([]byte(result), &cal); err == nil {
			return &cal, nil
		}
	}

	// 缓存未命中，从源获取
	cal, err := calendar.GetCalendar(dateStr)
	if err != nil {
		return nil, err
	}

	// 转换为 output 类型
	data := &output.CalendarData{
		Date:         cal.Date,
		LunarDate:    cal.LunarDate,
		Ganzhi:       cal.Ganzhi,
		Zodiac:       cal.Zodiac,
		DayOfYear:    cal.DayOfYear,
		SolarTerm:    cal.SolarTerm,
		Auspicious:   cal.Auspicious,
		Inauspicious: cal.Inauspicious,
	}

	// 缓存到 Redis，24 小时过期
	bytes, _ := json.Marshal(data)
	_ = u.redis.Set(ctx, cacheKey, string(bytes), 24*time.Hour)

	return data, nil
}

func (u *useCase) GetFooterLinks(ctx context.Context) ([]output.FooterLink, error) {
	links, err := u.footerLinks.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]output.FooterLink, len(links))
	for i, link := range links {
		result[i] = output.FooterLink{
			Title: link.Title,
			Link:  link.Link,
		}
	}
	return result, nil
}

func (u *useCase) loadWebsiteSettings(ctx context.Context) map[string]string {
	result := make(map[string]string)
	if u.settings == nil {
		return result
	}
	settings, err := u.settings.ListAll(ctx)
	if err != nil {
		return result
	}
	for _, item := range settings {
		result[item.SettingKey] = item.SettingValue
	}
	return result
}

func settingValueWithKey(values map[string]string, value string, keys ...string) output.SettingValueWithKey {
	for _, key := range keys {
		if values[key] != "" {
			return output.SettingValueWithKey{Value: value, SettingKey: key}
		}
	}
	if len(keys) > 0 {
		return output.SettingValueWithKey{Value: value, SettingKey: keys[0]}
	}
	return output.SettingValueWithKey{Value: value}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
