package file

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"

	"github.com/google/uuid"
)

var ErrRepo = errors.New("repo")

type useCase struct {
	files       repo.FileRepo
	objectStore repo.ObjectStore
}

// New 创建 File UseCase。
func New(files repo.FileRepo, objectStore repo.ObjectStore) usecase.File {
	return &useCase{
		files:       files,
		objectStore: objectStore,
	}
}

func (u *useCase) Upload(ctx context.Context, params input.UploadFile) (*output.UploadResult, error) {
	// 生成唯一的文件 key
	ext := filepath.Ext(params.Filename)
	key := fmt.Sprintf("%s/%s/%s%s",
		params.Usage,
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	// 上传到对象存储
	url, err := u.objectStore.Upload(ctx, key, params.File, params.Size, params.ContentType)
	if err != nil {
		return nil, fmt.Errorf("upload error: %w", err)
	}

	// 保存文件记录
	file := &entity.File{
		Key:      key,
		Filename: params.Filename,
		Size:     params.Size,
		MimeType: params.ContentType,
		Usage:    params.Usage,
	}

	_, err = u.files.Create(ctx, file)
	if err != nil {
		// 上传成功但保存记录失败，尝试删除已上传的文件
		_ = u.objectStore.Delete(ctx, key)
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return &output.UploadResult{
		Key: key,
		URL: url,
	}, nil
}

func (u *useCase) Delete(ctx context.Context, key string) error {
	// 删除对象存储中的文件
	if err := u.objectStore.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete object error: %w", err)
	}

	// 删除文件记录
	if err := u.files.Delete(ctx, key); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}

func (u *useCase) List(ctx context.Context, params input.ListFiles) (*output.ListResult[output.FileInfo], error) {
	offset := (params.Page - 1) * params.PageSize

	var filename, mimeType *string
	if params.Filename != "" {
		filename = &params.Filename
	}
	if params.MimeType != "" {
		mimeType = &params.MimeType
	}

	files, total, err := u.files.List(ctx, offset, params.PageSize, filename, mimeType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.FileInfo, len(files))
	for i, f := range files {
		items[i] = output.FileInfo{
			ID:        f.ID,
			Name:      f.Filename,
			URL:       u.objectStore.GetURL(f.Key),
			Size:      f.Size,
			MimeType:  f.MimeType,
			CreatedAt: f.CreatedAt,
		}
	}

	return &output.ListResult[output.FileInfo]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) DeleteByIDs(ctx context.Context, ids []int64) error {
	// 获取文件记录
	files, err := u.files.GetByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	// 删除对象存储中的文件
	for _, f := range files {
		if err := u.objectStore.Delete(ctx, f.Key); err != nil {
			// 记录错误但继续删除其他文件
			continue
		}
	}

	// 批量删除文件记录
	if err := u.files.DeleteByIDs(ctx, ids); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}
