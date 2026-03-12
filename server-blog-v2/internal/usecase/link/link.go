package link

import (
	"context"
	"errors"
	"fmt"

	"server-blog-v2/config"
	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
	"server-blog-v2/internal/usecase/urlutil"
)

var ErrRepo = errors.New("repo")

type useCase struct {
	cfg   *config.Config
	links repo.LinkRepo
}

// New 创建 Link UseCase。
func New(cfg *config.Config, links repo.LinkRepo) usecase.Link {
	return &useCase{cfg: cfg, links: links}
}

func (u *useCase) List(ctx context.Context) ([]*output.FriendLink, error) {
	links, err := u.links.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]*output.FriendLink, len(links))
	for i, l := range links {
		fl := &output.FriendLink{
			ID:   l.ID,
			Name: l.Name,
			Link: l.URL,
		}
		if l.Logo != nil {
			fl.Logo = urlutil.ResolveImageURL(u.cfg, *l.Logo)
		}
		if l.Description != nil {
			fl.Description = *l.Description
		}
		items[i] = fl
	}

	return items, nil
}

func (u *useCase) Create(ctx context.Context, params input.CreateLink) (int64, error) {
	link := &entity.Link{
		Name:      params.Name,
		URL:       params.URL,
		Sort:      params.Sort,
		IsVisible: true,
	}
	if params.Logo != "" {
		link.Logo = &params.Logo
	}
	if params.Description != "" {
		link.Description = &params.Description
	}

	id, err := u.links.Create(ctx, link)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return id, nil
}

func (u *useCase) Update(ctx context.Context, params input.UpdateLink) error {
	link := &entity.Link{
		ID:        params.ID,
		Name:      params.Name,
		URL:       params.URL,
		Sort:      params.Sort,
		IsVisible: true,
	}
	if params.Logo != "" {
		link.Logo = &params.Logo
	}
	if params.Description != "" {
		link.Description = &params.Description
	}

	if err := u.links.Update(ctx, link); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}

func (u *useCase) Delete(ctx context.Context, id int64) error {
	if err := u.links.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}
