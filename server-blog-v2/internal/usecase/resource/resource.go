package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"server-blog-v2/internal/entity"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
	"server-blog-v2/pkg/redis"
)

const (
	// ChunkSize 块大小 4MB
	ChunkSize = 4 * 1024 * 1024
	// TaskExpireHours 任务过期时间（小时）
	TaskExpireHours = 24 * 7 // 7天
	// DefaultMaxFileSize 默认最大文件大小 500MB
	DefaultMaxFileSize = 500 * 1024 * 1024
	// RedisKeyMaxFileSize Redis 中最大文件大小的 key
	RedisKeyMaxFileSize = "upload:max_size"
)

var ErrRepo = errors.New("repo")

type useCase struct {
	resources   repo.ResourceRepo
	tasks       repo.ResourceUploadTaskRepo
	objectStore repo.ObjectStore
	redis       redis.Client
}

// New 创建 Resource UseCase。
func New(resources repo.ResourceRepo, tasks repo.ResourceUploadTaskRepo, objectStore repo.ObjectStore, rdb redis.Client) usecase.Resource {
	return &useCase{
		resources:   resources,
		tasks:       tasks,
		objectStore: objectStore,
		redis:       rdb,
	}
}

func (u *useCase) List(ctx context.Context, userUUID *string, params input.ListResources) (*output.ListResult[output.ResourceInfo], error) {
	offset := (params.Page - 1) * params.PageSize

	var filename, mimeType *string
	if params.Filename != "" {
		filename = &params.Filename
	}
	if params.MimeType != "" {
		mimeType = &params.MimeType
	}

	resources, total, err := u.resources.List(ctx, offset, params.PageSize, userUUID, filename, mimeType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepo, err)
	}

	items := make([]output.ResourceInfo, len(resources))
	for i, r := range resources {
		items[i] = output.ResourceInfo{
			ID:              r.ID,
			FileKey:         r.FileKey,
			FileName:        r.FileName,
			FileURL:         u.objectStore.GetURL(r.FileKey),
			FileSize:        r.FileSize,
			MimeType:        r.MimeType,
			TranscodeStatus: int8(r.TranscodeStatus),
			CreatedAt:       r.CreatedAt,
		}
		if r.TranscodeKey != "" {
			url := u.objectStore.GetURL(r.TranscodeKey)
			items[i].TranscodeURL = &url
		}
		if r.ThumbnailKey != "" {
			url := u.objectStore.GetURL(r.ThumbnailKey)
			items[i].ThumbnailURL = &url
		}
	}

	return &output.ListResult[output.ResourceInfo]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) DeleteByIDs(ctx context.Context, ids []int64) error {
	// 获取资源记录
	resources, err := u.resources.GetByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	// 删除对象存储中的文件
	for _, r := range resources {
		_ = u.objectStore.Delete(ctx, r.FileKey)
		if r.TranscodeKey != "" {
			_ = u.objectStore.Delete(ctx, r.TranscodeKey)
		}
		if r.ThumbnailKey != "" {
			_ = u.objectStore.Delete(ctx, r.ThumbnailKey)
		}
	}

	// 批量删除资源记录
	if err := u.resources.DeleteByIDs(ctx, ids); err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	return nil
}

func (u *useCase) CheckFileHash(ctx context.Context, fileHash, userUUID string) (*entity.Resource, error) {
	resource, err := u.resources.GetByFileHash(ctx, fileHash, userUUID)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetMaxFileSize 获取最大文件大小。
// 优先从 Redis 的 "upload:max_size" 获取；若无值则使用默认值并写入 Redis。
func (u *useCase) GetMaxFileSize(ctx context.Context) int64 {
	val, err := u.redis.Get(ctx, RedisKeyMaxFileSize)
	if err == nil && val != "" {
		if size, err := strconv.ParseInt(val, 10, 64); err == nil {
			return size
		}
	}
	// Redis 中没有值（redis.Nil 或解析失败），写入默认值
	_ = u.redis.Set(ctx, RedisKeyMaxFileSize, strconv.FormatInt(DefaultMaxFileSize, 10), 0)
	return DefaultMaxFileSize
}

// Check 检查文件（秒传/续传检测）。
func (u *useCase) Check(ctx context.Context, userUUID string, params input.ResourceCheck) (*output.ResourceCheckResponse, error) {
	// 验证文件大小
	if params.FileSize > u.GetMaxFileSize(ctx) {
		return nil, errors.New("文件大小超过限制")
	}

	// 1. 检查当前用户是否已有相同 hash 的资源
	resource, err := u.resources.GetByFileHash(ctx, params.FileHash, userUUID)
	if err == nil && resource != nil {
		return &output.ResourceCheckResponse{
			Exists:  true,
			FileURL: u.objectStore.GetURL(resource.FileKey),
		}, nil
	}

	// 2. 检查其他用户是否已上传相同 hash 的资源（秒传）
	existingResource, err := u.resources.GetByFileHashAny(ctx, params.FileHash)
	if err == nil && existingResource != nil {
		// 为当前用户创建新记录指向同一物理文件
		newResource := &entity.Resource{
			FileKey:  existingResource.FileKey,
			FileName: params.FileName,
			FileHash: params.FileHash,
			FileSize: params.FileSize,
			MimeType: existingResource.MimeType,
			UserUUID: userUUID,
		}
		if _, err := u.resources.Create(ctx, newResource); err != nil {
			return nil, fmt.Errorf("创建资源记录失败: %w", err)
		}
		return &output.ResourceCheckResponse{
			Exists:  true,
			FileURL: u.objectStore.GetURL(existingResource.FileKey),
		}, nil
	}

	// 3. 检查是否存在未完成的任务（断点续传）
	existingTask, err := u.tasks.GetByFileHash(ctx, params.FileHash, userUUID, []entity.TaskStatus{entity.TaskStatusInit, entity.TaskStatusUploading})
	if err == nil && existingTask != nil {
		uploadedChunks, missingChunks := parseContexts(existingTask.QiniuContexts, existingTask.TotalChunks)
		return &output.ResourceCheckResponse{
			Exists:         false,
			TaskID:         existingTask.TaskID,
			TotalChunks:    existingTask.TotalChunks,
			UploadedChunks: uploadedChunks,
			MissingChunks:  missingChunks,
		}, nil
	}

	// 4. 新文件
	return &output.ResourceCheckResponse{Exists: false}, nil
}

// Init 初始化上传任务。
func (u *useCase) Init(ctx context.Context, userUUID string, params input.ResourceInit) (*output.ResourceInitResponse, error) {
	// 验证文件大小
	if params.FileSize > u.GetMaxFileSize(ctx) {
		return nil, errors.New("文件大小超过限制")
	}

	// 计算总块数
	totalChunks := int(math.Ceil(float64(params.FileSize) / float64(ChunkSize)))
	if totalChunks == 0 {
		totalChunks = 1
	}

	// 初始化 QiniuContexts 数组
	contexts := make([]string, totalChunks)
	contextsJSON, _ := json.Marshal(contexts)

	// 生成任务 ID
	taskUUID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("生成任务ID失败: %w", err)
	}

	// 创建任务
	task := &entity.ResourceUploadTask{
		TaskID:        taskUUID.String(),
		FileName:      params.FileName,
		FileSize:      params.FileSize,
		FileHash:      params.FileHash,
		MimeType:      params.MimeType,
		ChunkSize:     ChunkSize,
		TotalChunks:   totalChunks,
		Status:        entity.TaskStatusInit,
		UserUUID:      userUUID,
		ExpiresAt:     time.Now().Add(TaskExpireHours * time.Hour),
		QiniuContexts: string(contextsJSON),
	}

	if _, err := u.tasks.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	return &output.ResourceInitResponse{
		TaskID:      task.TaskID,
		TotalChunks: totalChunks,
		ChunkSize:   ChunkSize,
	}, nil
}

// UploadChunk 上传分片。
func (u *useCase) UploadChunk(ctx context.Context, userUUID string, params input.ResourceUploadChunk) (*output.ResourceUploadChunkResponse, error) {
	// 查询任务
	task, err := u.tasks.GetByTaskID(ctx, params.TaskID, userUUID)
	if err != nil {
		return nil, errors.New("任务不存在")
	}

	// 检查任务状态
	if task.Status == entity.TaskStatusCompleted {
		return nil, errors.New("任务已完成")
	}
	if task.Status == entity.TaskStatusCancelled {
		return nil, errors.New("任务已取消")
	}

	// 检查块号有效性
	if params.ChunkNumber < 0 || params.ChunkNumber >= task.TotalChunks {
		return nil, fmt.Errorf("无效的块号: %d, 有效范围: 0-%d", params.ChunkNumber, task.TotalChunks-1)
	}

	// 流式上传块到七牛云
	qiniuCtx, err := u.objectStore.UploadBlock(ctx, params.ChunkData, params.ChunkSize)
	if err != nil {
		return nil, fmt.Errorf("上传块失败: %w", err)
	}

	// 更新任务
	if err := u.tasks.UpdateChunkContext(ctx, params.TaskID, params.ChunkNumber, qiniuCtx, entity.TaskStatusUploading); err != nil {
		return nil, fmt.Errorf("更新任务失败: %w", err)
	}

	return &output.ResourceUploadChunkResponse{
		Success:     true,
		ChunkNumber: params.ChunkNumber,
	}, nil
}

// Complete 完成上传。
func (u *useCase) Complete(ctx context.Context, userUUID string, params input.ResourceComplete) (*output.ResourceCompleteResponse, error) {
	// 查询任务
	task, err := u.tasks.GetByTaskID(ctx, params.TaskID, userUUID)
	if err != nil {
		return nil, errors.New("任务不存在")
	}

	// 检查任务状态
	if task.Status == entity.TaskStatusCompleted {
		return nil, errors.New("任务已完成")
	}

	// 解析 contexts
	var contexts []string
	if err := json.Unmarshal([]byte(task.QiniuContexts), &contexts); err != nil {
		return nil, fmt.Errorf("解析contexts失败: %w", err)
	}

	// 检查是否所有块都已上传
	var validContexts []string
	for i, ctx := range contexts {
		if ctx == "" {
			return nil, fmt.Errorf("块 %d 尚未上传", i)
		}
		validContexts = append(validContexts, ctx)
	}

	// 生成文件 Key
	fileKey := u.objectStore.GenerateFileKey(task.FileName, task.FileHash)

	// 合并文件
	if err := u.objectStore.MergeBlocks(ctx, task.FileSize, fileKey, validContexts); err != nil {
		_ = u.tasks.UpdateStatus(ctx, task.TaskID, entity.TaskStatusFailed)
		return nil, fmt.Errorf("合并文件失败: %w", err)
	}

	// 判断是否需要转码（视频文件）
	transcodeStatus := entity.TranscodeStatusNone
	if isVideoMimeType(task.MimeType) {
		transcodeStatus = entity.TranscodeStatusProcessing
	}

	// 创建资源记录
	resourceRecord := &entity.Resource{
		FileKey:         fileKey,
		FileName:        task.FileName,
		FileHash:        task.FileHash,
		FileSize:        task.FileSize,
		MimeType:        task.MimeType,
		UserUUID:        userUUID,
		TranscodeStatus: transcodeStatus,
	}
	if _, err := u.resources.Create(ctx, resourceRecord); err != nil {
		return nil, fmt.Errorf("创建资源记录失败: %w", err)
	}

	// 更新任务状态
	_ = u.tasks.UpdateStatus(ctx, task.TaskID, entity.TaskStatusCompleted)

	return &output.ResourceCompleteResponse{
		FileURL: u.objectStore.GetURL(fileKey),
		FileKey: fileKey,
	}, nil
}

// Cancel 取消上传。
func (u *useCase) Cancel(ctx context.Context, userUUID string, params input.ResourceCancel) error {
	task, err := u.tasks.GetByTaskID(ctx, params.TaskID, userUUID)
	if err != nil {
		return errors.New("任务不存在")
	}

	if task.Status == entity.TaskStatusCompleted {
		return errors.New("任务已完成，无法取消")
	}

	return u.tasks.UpdateStatus(ctx, params.TaskID, entity.TaskStatusCancelled)
}

// Progress 查询上传进度。
func (u *useCase) Progress(ctx context.Context, userUUID string, params input.ResourceProgress) (*output.ResourceProgressResponse, error) {
	task, err := u.tasks.GetByTaskID(ctx, params.TaskID, userUUID)
	if err != nil {
		return nil, errors.New("任务不存在")
	}

	uploadedChunks, missingChunks := parseContexts(task.QiniuContexts, task.TotalChunks)
	progress := 0
	if task.TotalChunks > 0 {
		progress = len(uploadedChunks) * 100 / task.TotalChunks
	}

	return &output.ResourceProgressResponse{
		TaskID:         task.TaskID,
		TotalChunks:    task.TotalChunks,
		UploadedChunks: uploadedChunks,
		MissingChunks:  missingChunks,
		Progress:       progress,
	}, nil
}

// HandleQiniuCallback 处理七牛云转码回调。
func (u *useCase) HandleQiniuCallback(ctx context.Context, inputKey string, code int, items []input.QiniuCallbackItem) error {
	if code != 0 {
		return u.resources.UpdateTranscodeStatusByFileKey(ctx, inputKey, entity.TranscodeStatusFailed, "", "")
	}

	var transcodeKey, thumbnailKey string
	for _, item := range items {
		if item.Code != 0 {
			continue
		}
		if strings.HasSuffix(item.Key, "_h264.mp4") {
			transcodeKey = item.Key
		} else if strings.HasSuffix(item.Key, "_thumb.jpg") {
			thumbnailKey = item.Key
		}
	}

	return u.resources.UpdateTranscodeStatusByFileKey(ctx, inputKey, entity.TranscodeStatusSuccess, transcodeKey, thumbnailKey)
}

// parseContexts 解析 QiniuContexts，返回已上传和缺失的块号。
func parseContexts(contextsJSON string, totalChunks int) ([]int, []int) {
	var contexts []string
	if err := json.Unmarshal([]byte(contextsJSON), &contexts); err != nil {
		return nil, nil
	}

	var uploadedChunks, missingChunks []int
	for i := 0; i < totalChunks; i++ {
		if i < len(contexts) && contexts[i] != "" {
			uploadedChunks = append(uploadedChunks, i)
		} else {
			missingChunks = append(missingChunks, i)
		}
	}

	return uploadedChunks, missingChunks
}

// isVideoMimeType 判断是否为视频 MIME 类型。
func isVideoMimeType(mimeType string) bool {
	return len(mimeType) >= 6 && mimeType[:6] == "video/"
}
