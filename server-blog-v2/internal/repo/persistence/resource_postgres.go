package persistence

import (
	"context"
	"time"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"

	"gorm.io/gorm"
)

type resourceRepo struct {
	db *gorm.DB
}

// NewResourceRepo 创建资源仓库。
func NewResourceRepo(db *gorm.DB) repo.ResourceRepo {
	return &resourceRepo{db: db}
}

func (r *resourceRepo) Create(ctx context.Context, resource *entity.Resource) (int64, error) {
	mr := toModelResource(resource)
	if err := r.db.WithContext(ctx).Create(mr).Error; err != nil {
		return 0, err
	}
	return mr.ID, nil
}

func (r *resourceRepo) GetByID(ctx context.Context, id int64) (*entity.Resource, error) {
	var mr model.Resource
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&mr).Error; err != nil {
		return nil, err
	}
	return toEntityResource(&mr), nil
}

func (r *resourceRepo) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Resource, error) {
	var mrs []model.Resource
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&mrs).Error; err != nil {
		return nil, err
	}

	resources := make([]*entity.Resource, len(mrs))
	for i, mr := range mrs {
		resources[i] = toEntityResource(&mr)
	}
	return resources, nil
}

func (r *resourceRepo) GetByFileHash(ctx context.Context, fileHash, userUUID string) (*entity.Resource, error) {
	var mr model.Resource
	if err := r.db.WithContext(ctx).Where("file_hash = ? AND user_uuid = ?", fileHash, userUUID).First(&mr).Error; err != nil {
		return nil, err
	}
	return toEntityResource(&mr), nil
}

func (r *resourceRepo) GetByFileHashAny(ctx context.Context, fileHash string) (*entity.Resource, error) {
	var mr model.Resource
	if err := r.db.WithContext(ctx).Where("file_hash = ?", fileHash).First(&mr).Error; err != nil {
		return nil, err
	}
	return toEntityResource(&mr), nil
}

func (r *resourceRepo) List(ctx context.Context, offset, limit int, userUUID *string, filename, mimeType *string) ([]*entity.Resource, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Resource{})

	if userUUID != nil && *userUUID != "" {
		db = db.Where("user_uuid = ?", *userUUID)
	}
	if filename != nil && *filename != "" {
		db = db.Where("file_name LIKE ?", "%"+*filename+"%")
	}
	if mimeType != nil && *mimeType != "" {
		db = db.Where("mime_type LIKE ?", *mimeType+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var mrs []model.Resource
	if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&mrs).Error; err != nil {
		return nil, 0, err
	}

	resources := make([]*entity.Resource, len(mrs))
	for i, mr := range mrs {
		resources[i] = toEntityResource(&mr)
	}
	return resources, total, nil
}

func (r *resourceRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Resource{}).Error
}

func (r *resourceRepo) DeleteByIDs(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Resource{}).Error
}

func (r *resourceRepo) UpdateTranscodeStatus(ctx context.Context, id int64, status entity.TranscodeStatus, transcodeKey, thumbnailKey string) error {
	updates := map[string]interface{}{
		"transcode_status": int16(status),
		"updated_at":       time.Now(),
	}
	if transcodeKey != "" {
		updates["transcode_key"] = transcodeKey
	}
	if thumbnailKey != "" {
		updates["thumbnail_key"] = thumbnailKey
	}
	return r.db.WithContext(ctx).Model(&model.Resource{}).Where("id = ?", id).Updates(updates).Error
}

func (r *resourceRepo) UpdateTranscodeStatusByFileKey(ctx context.Context, fileKey string, status entity.TranscodeStatus, transcodeKey, thumbnailKey string) error {
	updates := map[string]interface{}{
		"transcode_status": int16(status),
		"updated_at":       time.Now(),
	}
	if transcodeKey != "" {
		updates["transcode_key"] = transcodeKey
	}
	if thumbnailKey != "" {
		updates["thumbnail_key"] = thumbnailKey
	}
	return r.db.WithContext(ctx).Model(&model.Resource{}).Where("file_key = ?", fileKey).Updates(updates).Error
}

func toModelResource(r *entity.Resource) *model.Resource {
	mr := &model.Resource{
		ID:       r.ID,
		FileKey:  r.FileKey,
		FileName: r.FileName,
		FileSize: r.FileSize,
		MimeType: r.MimeType,
	}
	if r.FileHash != "" {
		mr.FileHash = &r.FileHash
	}
	if r.UserUUID != "" {
		mr.UserUUID = &r.UserUUID
	}
	status := int16(r.TranscodeStatus)
	mr.TranscodeStatus = &status
	if r.TranscodeKey != "" {
		mr.TranscodeKey = &r.TranscodeKey
	}
	if r.ThumbnailKey != "" {
		mr.ThumbnailKey = &r.ThumbnailKey
	}
	return mr
}

func toEntityResource(mr *model.Resource) *entity.Resource {
	r := &entity.Resource{
		ID:       mr.ID,
		FileKey:  mr.FileKey,
		FileName: mr.FileName,
		FileSize: mr.FileSize,
		MimeType: mr.MimeType,
	}
	if mr.FileHash != nil {
		r.FileHash = *mr.FileHash
	}
	if mr.UserUUID != nil {
		r.UserUUID = *mr.UserUUID
	}
	if mr.TranscodeStatus != nil {
		r.TranscodeStatus = entity.TranscodeStatus(*mr.TranscodeStatus)
	}
	if mr.TranscodeKey != nil {
		r.TranscodeKey = *mr.TranscodeKey
	}
	if mr.ThumbnailKey != nil {
		r.ThumbnailKey = *mr.ThumbnailKey
	}
	if mr.CreatedAt != nil {
		r.CreatedAt = *mr.CreatedAt
	}
	if mr.UpdatedAt != nil {
		r.UpdatedAt = *mr.UpdatedAt
	}
	return r
}
