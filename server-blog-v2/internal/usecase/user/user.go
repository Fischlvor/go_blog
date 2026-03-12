package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"server-blog-v2/config"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
	"server-blog-v2/internal/usecase/urlutil"
)

var (
	ErrRepo     = errors.New("repo")
	ErrNotFound = errors.New("not found")
)

type useCase struct {
	cfg   *config.Config
	users repo.UserRepo
}

// New 创建 User UseCase。
func New(cfg *config.Config, users repo.UserRepo) usecase.User {
	return &useCase{cfg: cfg, users: users}
}

func (u *useCase) GetProfile(ctx context.Context, userUUID string) (*output.UserProfile, error) {
	user, err := u.users.GetByUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	if user == nil {
		return nil, ErrNotFound
	}

	profile := &output.UserProfile{
		ID:        user.ID,
		UUID:      user.UUID,
		Nickname:  user.Nickname,
		Avatar:    urlutil.ResolveImageURL(u.cfg, user.Avatar),
		Signature: "签名是空白的，这位用户似乎比较低调。",
		RoleID:    user.RoleID,
	}
	if user.Email != nil {
		profile.Email = *user.Email
	}

	return profile, nil
}

func (u *useCase) UpdateProfile(ctx context.Context, userUUID string, params input.UpdateProfile) error {
	user, err := u.users.GetByUUID(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	if user == nil {
		return ErrNotFound
	}

	// 更新字段
	if params.Nickname != "" {
		user.Nickname = params.Nickname
	}
	if params.Avatar != "" {
		user.Avatar = params.Avatar
	}

	if err := u.users.Update(ctx, user); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}

// List 用户列表（管理端）。
func (u *useCase) List(ctx context.Context, params input.ListUsers) (*output.ListResult[output.UserAdmin], error) {
	offset := (params.Page - 1) * params.PageSize

	var keyword *string
	if params.Keyword != nil {
		keyword = &params.Keyword.Keyword
	}

	users, total, err := u.users.List(ctx, offset, params.PageSize, keyword, params.Status)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.UserAdmin, len(users))
	for i, user := range users {
		items[i] = output.UserAdmin{
			ID:        user.ID,
			UUID:      user.UUID,
			Username:  user.Nickname,
			Email:     "",
			Avatar:    urlutil.ResolveImageURL(u.cfg, user.Avatar),
			RoleID:    user.RoleID,
			Freeze:    user.Status == "frozen" || user.Status == "disabled",
			CreatedAt: user.CreatedAt,
		}
		if user.Email != nil {
			items[i].Email = *user.Email
		}
	}

	return &output.ListResult[output.UserAdmin]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

// Freeze 冻结用户。
func (u *useCase) Freeze(ctx context.Context, userUUID string) error {
	if err := u.users.UpdateStatus(ctx, userUUID, "frozen"); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

// Unfreeze 解冻用户。
func (u *useCase) Unfreeze(ctx context.Context, userUUID string) error {
	if err := u.users.UpdateStatus(ctx, userUUID, "active"); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

// GetChart 获取用户图表数据。
func (u *useCase) GetChart(ctx context.Context, days int) (*output.UserChart, error) {
	result := &output.UserChart{
		DateList:     make([]string, days),
		LoginData:    make([]int, days),
		RegisterData: make([]int, days),
	}

	// 生成日期列表
	startDate := time.Now().AddDate(0, 0, -days)
	for i := 0; i < days; i++ {
		result.DateList[i] = startDate.AddDate(0, 0, i+1).Format("2006-01-02")
	}

	// 获取登录和注册数据
	loginCounts, err := u.users.GetLoginCountsByDate(ctx, days)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	registerCounts, err := u.users.GetRegisterCountsByDate(ctx, days)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	// 填充数据
	for i, date := range result.DateList {
		if count, ok := loginCounts[date]; ok {
			result.LoginData[i] = count
		}
		if count, ok := registerCounts[date]; ok {
			result.RegisterData[i] = count
		}
	}

	return result, nil
}

// ListLogins 登录记录列表。
func (u *useCase) ListLogins(ctx context.Context, params input.ListLogins) (*output.ListResult[output.LoginInfo], error) {
	offset := (params.Page - 1) * params.PageSize

	var keyword *string
	if params.Keyword != nil {
		keyword = &params.Keyword.Keyword
	}

	logins, total, err := u.users.ListLogins(ctx, offset, params.PageSize, keyword)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.LoginInfo, len(logins))
	for i, login := range logins {
		items[i] = output.LoginInfo{
			ID:          login.ID,
			UserUUID:    login.UserUUID,
			LoginMethod: login.LoginMethod,
			IP:          login.IP,
			Address:     login.Address,
			OS:          login.OS,
			DeviceInfo:  login.DeviceInfo,
			BrowserInfo: login.BrowserInfo,
			Status:      login.Status,
			CreatedAt:   login.CreatedAt,
		}
	}

	return &output.ListResult[output.LoginInfo]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}
