package persistence

import (
	"context"
	"time"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type userRepo struct {
	query *query.Query
}

// NewUserRepo 创建用户仓库。
func NewUserRepo(db *gorm.DB) repo.UserRepo {
	return &userRepo{query: query.Use(db)}
}

func (r *userRepo) GetByUUID(ctx context.Context, uuid string) (*entity.User, error) {
	u := r.query.User
	row, err := u.WithContext(ctx).Where(u.UUID.Eq(uuid)).First()
	if err != nil {
		return nil, err
	}
	return toEntityUser(row), nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	u := r.query.User
	row, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toEntityUser(row), nil
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) (int64, error) {
	mu := toModelUser(user)
	if err := r.query.User.WithContext(ctx).Create(mu); err != nil {
		return 0, err
	}
	return mu.ID, nil
}

func (r *userRepo) CreateFromSSO(ctx context.Context, uuid, nickname, email, avatar, address, signature string, registerSource int) error {
	now := time.Now()
	roleID := int32(1) // 默认普通用户
	mu := &model.User{
		UUID:      uuid,
		Nickname:  &nickname,
		Email:     &email,
		Avatar:    &avatar,
		RoleID:    &roleID,
		Status:    ptrString("active"),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	return r.query.User.WithContext(ctx).Create(mu)
}

func ptrString(s string) *string {
	return &s
}

func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	u := r.query.User
	mu := toModelUser(user)
	_, err := u.WithContext(ctx).Where(u.ID.Eq(user.ID)).Updates(mu)
	return err
}

func (r *userRepo) List(ctx context.Context, offset, limit int, keyword, status *string) ([]*entity.User, int64, error) {
	u := r.query.User
	do := u.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		do = do.Where(u.Nickname.Like("%" + *keyword + "%"))
	}
	if status != nil && *status != "" {
		do = do.Where(u.Status.Eq(*status))
	}

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(u.CreatedAt.Desc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(rows))
	for i, row := range rows {
		users[i] = toEntityUser(row)
	}
	return users, total, nil
}

func (r *userRepo) UpdateStatus(ctx context.Context, uuid string, status string) error {
	u := r.query.User
	_, err := u.WithContext(ctx).Where(u.UUID.Eq(uuid)).UpdateSimple(u.Status.Value(status))
	return err
}

func (r *userRepo) GetRoleByUUID(ctx context.Context, uuid string) (int, error) {
	u := r.query.User
	row, err := u.WithContext(ctx).Where(u.UUID.Eq(uuid)).First()
	if err != nil {
		return 0, err
	}
	if row.RoleID != nil {
		return int(*row.RoleID), nil
	}
	return 0, nil
}

func toModelUser(u *entity.User) *model.User {
	roleID := int32(u.RoleID)
	mu := &model.User{
		ID:       u.ID,
		UUID:     u.UUID,
		Nickname: &u.Nickname,
		Avatar:   &u.Avatar,
		RoleID:   &roleID,
		Status:   &u.Status,
	}
	if u.Email != nil {
		mu.Email = u.Email
	}
	return mu
}

func toEntityUser(mu *model.User) *entity.User {
	user := &entity.User{
		ID:    mu.ID,
		UUID:  mu.UUID,
		Email: mu.Email,
	}
	if mu.Nickname != nil {
		user.Nickname = *mu.Nickname
	}
	if mu.Avatar != nil {
		user.Avatar = *mu.Avatar
	}
	if mu.RoleID != nil {
		user.RoleID = int(*mu.RoleID)
	}
	if mu.Status != nil {
		user.Status = *mu.Status
	}
	if mu.CreatedAt != nil {
		user.CreatedAt = *mu.CreatedAt
	}
	if mu.UpdatedAt != nil {
		user.UpdatedAt = *mu.UpdatedAt
	}
	return user
}

func (r *userRepo) GetLoginCountsByDate(ctx context.Context, days int) (map[string]int, error) {
	// 由于 v2 使用 SSO 登录，登录记录在 SSO 服务中
	// 这里返回空数据，后续可以对接 SSO 服务的登录统计
	return make(map[string]int), nil
}

func (r *userRepo) GetRegisterCountsByDate(ctx context.Context, days int) (map[string]int, error) {
	u := r.query.User
	startDate := time.Now().AddDate(0, 0, -days)

	// 查询最近 N 天注册的用户
	rows, err := u.WithContext(ctx).
		Where(u.CreatedAt.Gte(startDate)).
		Find()
	if err != nil {
		return nil, err
	}

	// 按日期统计
	counts := make(map[string]int)
	for _, row := range rows {
		if row.CreatedAt != nil {
			date := row.CreatedAt.Format("2006-01-02")
			counts[date]++
		}
	}

	return counts, nil
}

func (r *userRepo) ListLogins(ctx context.Context, offset, limit int, keyword *string) ([]*entity.Login, int64, error) {
	l := r.query.Login
	q := l.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		q = q.Where(l.UserUUID.Like("%" + *keyword + "%")).
			Or(l.IP.Like("%" + *keyword + "%")).
			Or(l.Address.Like("%" + *keyword + "%"))
	}

	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.Order(l.CreatedAt.Desc()).Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	logins := make([]*entity.Login, len(rows))
	for i, row := range rows {
		logins[i] = toEntityLogin(row)
	}

	return logins, total, nil
}

func toEntityLogin(m *model.Login) *entity.Login {
	login := &entity.Login{
		ID: m.ID,
	}
	if m.UserUUID != nil {
		login.UserUUID = *m.UserUUID
	}
	if m.LoginMethod != nil {
		login.LoginMethod = *m.LoginMethod
	}
	if m.IP != nil {
		login.IP = *m.IP
	}
	if m.Address != nil {
		login.Address = *m.Address
	}
	if m.Os != nil {
		login.OS = *m.Os
	}
	if m.DeviceInfo != nil {
		login.DeviceInfo = *m.DeviceInfo
	}
	if m.BrowserInfo != nil {
		login.BrowserInfo = *m.BrowserInfo
	}
	if m.Status != nil {
		login.Status = int(*m.Status)
	}
	if m.CreatedAt != nil {
		login.CreatedAt = *m.CreatedAt
	}
	if m.UpdatedAt != nil {
		login.UpdatedAt = *m.UpdatedAt
	}
	return login
}
