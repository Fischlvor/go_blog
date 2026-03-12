package persistence

import (
	"context"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/repo/persistence/gen/model"
	"server-blog-v2/internal/repo/persistence/gen/query"

	"gorm.io/gorm"
)

type fileRepo struct {
	query *query.Query
}

// NewFileRepo 创建文件仓库。
func NewFileRepo(db *gorm.DB) repo.FileRepo {
	return &fileRepo{query: query.Use(db)}
}

func (r *fileRepo) Create(ctx context.Context, file *entity.File) (int64, error) {
	mf := toModelFile(file)
	if err := r.query.File.WithContext(ctx).Create(mf); err != nil {
		return 0, err
	}
	return mf.ID, nil
}

func (r *fileRepo) GetByKey(ctx context.Context, key string) (*entity.File, error) {
	f := r.query.File
	row, err := f.WithContext(ctx).Where(f.Key.Eq(key)).First()
	if err != nil {
		return nil, err
	}
	return toEntityFile(row), nil
}

func (r *fileRepo) List(ctx context.Context, offset, limit int, filename, mimeType *string) ([]*entity.File, int64, error) {
	f := r.query.File
	do := f.WithContext(ctx)

	if filename != nil && *filename != "" {
		do = do.Where(f.Filename.Like("%" + *filename + "%"))
	}
	if mimeType != nil && *mimeType != "" {
		do = do.Where(f.MimeType.Like(*mimeType + "%"))
	}

	total, err := do.Count()
	if err != nil {
		return nil, 0, err
	}

	do = do.Order(f.CreatedAt.Desc())
	rows, err := do.Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, 0, err
	}

	files := make([]*entity.File, len(rows))
	for i, row := range rows {
		files[i] = toEntityFile(row)
	}
	return files, total, nil
}

func (r *fileRepo) GetByIDs(ctx context.Context, ids []int64) ([]*entity.File, error) {
	f := r.query.File
	rows, err := f.WithContext(ctx).Where(f.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	files := make([]*entity.File, len(rows))
	for i, row := range rows {
		files[i] = toEntityFile(row)
	}
	return files, nil
}

func (r *fileRepo) Delete(ctx context.Context, key string) error {
	f := r.query.File
	_, err := f.WithContext(ctx).Where(f.Key.Eq(key)).Delete()
	return err
}

func (r *fileRepo) DeleteByIDs(ctx context.Context, ids []int64) error {
	f := r.query.File
	_, err := f.WithContext(ctx).Where(f.ID.In(ids...)).Delete()
	return err
}

func toModelFile(f *entity.File) *model.File {
	mf := &model.File{
		ID:       f.ID,
		Key:      f.Key,
		Filename: &f.Filename,
		Size:     &f.Size,
		MimeType: &f.MimeType,
		Usage:    &f.Usage,
	}
	if f.ResourceID != nil {
		mf.ResourceID = f.ResourceID
	}
	return mf
}

func toEntityFile(mf *model.File) *entity.File {
	file := &entity.File{
		ID:         mf.ID,
		Key:        mf.Key,
		ResourceID: mf.ResourceID,
	}
	if mf.Filename != nil {
		file.Filename = *mf.Filename
	}
	if mf.Size != nil {
		file.Size = *mf.Size
	}
	if mf.MimeType != nil {
		file.MimeType = *mf.MimeType
	}
	if mf.Usage != nil {
		file.Usage = *mf.Usage
	}
	if mf.CreatedAt != nil {
		file.CreatedAt = *mf.CreatedAt
	}
	return file
}
