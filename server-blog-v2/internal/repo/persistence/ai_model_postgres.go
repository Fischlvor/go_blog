package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"

	"gorm.io/gorm"
)

type aiModelRepo struct {
	db *gorm.DB
}

// NewAIModelRepo 创建 AI 模型仓库。
func NewAIModelRepo(db *gorm.DB) repo.AIModelRepo {
	return &aiModelRepo{db: db}
}

func (r *aiModelRepo) Create(ctx context.Context, m *entity.AIModel) (int64, error) {
	mm := toModelAIModel(m)
	if err := r.db.WithContext(ctx).Create(mm).Error; err != nil {
		return 0, err
	}
	return mm.ID, nil
}

func (r *aiModelRepo) GetByID(ctx context.Context, id int64) (*entity.AIModel, error) {
	var mm model.AiModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&mm).Error; err != nil {
		return nil, err
	}
	return toEntityAIModel(&mm), nil
}

func (r *aiModelRepo) List(ctx context.Context, offset, limit int, name, provider *string) ([]*entity.AIModel, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.AiModel{})

	if name != nil && *name != "" {
		db = db.Where("name LIKE ? OR display_name LIKE ?", "%"+*name+"%", "%"+*name+"%")
	}
	if provider != nil && *provider != "" {
		db = db.Where("provider = ?", *provider)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var mms []model.AiModel
	if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&mms).Error; err != nil {
		return nil, 0, err
	}

	models := make([]*entity.AIModel, len(mms))
	for i, mm := range mms {
		models[i] = toEntityAIModel(&mm)
	}
	return models, total, nil
}

func (r *aiModelRepo) Update(ctx context.Context, m *entity.AIModel) error {
	mm := toModelAIModel(m)
	return r.db.WithContext(ctx).Save(mm).Error
}

func (r *aiModelRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.AiModel{}).Error
}

func toModelAIModel(m *entity.AIModel) *model.AiModel {
	mm := &model.AiModel{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Provider:    m.Provider,
	}
	if m.Endpoint != "" {
		mm.Endpoint = &m.Endpoint
	}
	if m.ApiKey != "" {
		mm.APIKey = &m.ApiKey
	}
	maxTokens := int32(m.MaxTokens)
	mm.MaxTokens = &maxTokens
	mm.Temperature = &m.Temperature
	mm.IsActive = &m.IsActive
	return mm
}

func toEntityAIModel(mm *model.AiModel) *entity.AIModel {
	m := &entity.AIModel{
		ID:          mm.ID,
		Name:        mm.Name,
		DisplayName: mm.DisplayName,
		Provider:    mm.Provider,
	}
	if mm.Endpoint != nil {
		m.Endpoint = *mm.Endpoint
	}
	if mm.APIKey != nil {
		m.ApiKey = *mm.APIKey
	}
	if mm.MaxTokens != nil {
		m.MaxTokens = int(*mm.MaxTokens)
	}
	if mm.Temperature != nil {
		m.Temperature = *mm.Temperature
	}
	if mm.IsActive != nil {
		m.IsActive = *mm.IsActive
	}
	if mm.CreatedAt != nil {
		m.CreatedAt = *mm.CreatedAt
	}
	if mm.UpdatedAt != nil {
		m.UpdatedAt = *mm.UpdatedAt
	}
	return m
}
