package aimodel

import (
	"context"
	"errors"
	"fmt"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
)

var ErrRepo = errors.New("repo")

type useCase struct {
	models repo.AIModelRepo
}

// New 创建 AIModel UseCase。
func New(models repo.AIModelRepo) usecase.AIModel {
	return &useCase{models: models}
}

func (u *useCase) List(ctx context.Context, params input.ListAIModels) (*output.ListResult[output.AIModelInfo], error) {
	offset := (params.Page - 1) * params.PageSize

	var name, provider *string
	if params.Name != "" {
		name = &params.Name
	}
	if params.Provider != "" {
		provider = &params.Provider
	}

	models, total, err := u.models.List(ctx, offset, params.PageSize, name, provider)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.AIModelInfo, len(models))
	for i, m := range models {
		items[i] = toOutput(m)
	}

	return &output.ListResult[output.AIModelInfo]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) GetByID(ctx context.Context, id int64) (*output.AIModelInfo, error) {
	m, err := u.models.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	info := toOutput(m)
	return &info, nil
}

func (u *useCase) Create(ctx context.Context, params input.CreateAIModel) (int64, error) {
	m := &entity.AIModel{
		Name:        params.Name,
		DisplayName: params.DisplayName,
		Provider:    params.Provider,
		Endpoint:    params.Endpoint,
		ApiKey:      params.ApiKey,
		MaxTokens:   params.MaxTokens,
		Temperature: params.Temperature,
		IsActive:    params.IsActive,
	}
	id, err := u.models.Create(ctx, m)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return id, nil
}

func (u *useCase) Update(ctx context.Context, params input.UpdateAIModel) error {
	m := &entity.AIModel{
		ID:          params.ID,
		Name:        params.Name,
		DisplayName: params.DisplayName,
		Provider:    params.Provider,
		Endpoint:    params.Endpoint,
		ApiKey:      params.ApiKey,
		MaxTokens:   params.MaxTokens,
		Temperature: params.Temperature,
		IsActive:    params.IsActive,
	}
	if err := u.models.Update(ctx, m); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

func (u *useCase) Delete(ctx context.Context, id int64) error {
	if err := u.models.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}
	return nil
}

func toOutput(m *entity.AIModel) output.AIModelInfo {
	return output.AIModelInfo{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Provider:    m.Provider,
		Endpoint:    m.Endpoint,
		MaxTokens:   m.MaxTokens,
		Temperature: m.Temperature,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
