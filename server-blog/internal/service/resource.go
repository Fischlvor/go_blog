package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"server/internal/model/database"
	"server/internal/model/request"
	"server/internal/model/response"
	"server/pkg/global"
	"server/pkg/resource"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

const (
	// ChunkSize 块大小 4MB
	ChunkSize = 4 * 1024 * 1024
	// TaskExpireHours 任务过期时间（小时）
	TaskExpireHours = 24 * 7 // 7天
)

type ResourceService struct{}

// getUploader 获取上传器（延迟初始化）
func (s *ResourceService) getUploader() *resource.QiniuResourceUploader {
	return resource.NewQiniuResourceUploader()
}

// getValidator 获取验证器
func (s *ResourceService) getValidator() *resource.DefaultValidator {
	return resource.NewDefaultValidator()
}

// Check 检查文件（秒传/续传检测）
func (s *ResourceService) Check(req request.ResourceCheck, userUUID uuid.UUID) (*response.ResourceCheckResponse, error) {
	// 1. 检查当前用户是否已有相同hash的资源（避免重复记录）
	var userResource database.Resource
	err := global.DB.Where("file_hash = ? AND user_uuid = ?", req.FileHash, userUUID).First(&userResource).Error
	if err == nil {
		// 当前用户已有该文件，直接返回
		return &response.ResourceCheckResponse{
			Exists:  true,
			FileURL: s.getUploader().GetFileURL(userResource.FileKey),
		}, nil
	}

	// 2. 检查其他用户是否已上传相同hash的资源（秒传）
	var existingResource database.Resource
	err = global.DB.Where("file_hash = ?", req.FileHash).First(&existingResource).Error
	if err == nil {
		// 找到相同hash的资源，为当前用户创建新记录指向同一物理文件
		newResource := database.Resource{
			FileKey:  existingResource.FileKey,
			FileName: req.FileName,
			FileHash: req.FileHash,
			FileSize: req.FileSize,
			MimeType: existingResource.MimeType,
			UserUUID: userUUID,
		}
		if err := global.DB.Create(&newResource).Error; err != nil {
			return nil, fmt.Errorf("创建资源记录失败: %w", err)
		}

		return &response.ResourceCheckResponse{
			Exists:  true,
			FileURL: s.getUploader().GetFileURL(existingResource.FileKey),
		}, nil
	}

	// 2. 检查是否存在未完成的任务（断点续传）
	var existingTask database.ResourceUploadTask
	err = global.DB.Where("file_hash = ? AND user_uuid = ? AND status IN (?, ?)",
		req.FileHash, userUUID, database.TaskStatusInit, database.TaskStatusUploading).
		First(&existingTask).Error
	if err == nil {
		// 找到未完成的任务，返回任务信息
		uploadedChunks, missingChunks := s.parseContexts(existingTask.QiniuContexts, existingTask.TotalChunks)
		return &response.ResourceCheckResponse{
			Exists:         false,
			TaskID:         existingTask.TaskID,
			TotalChunks:    existingTask.TotalChunks,
			UploadedChunks: uploadedChunks,
			MissingChunks:  missingChunks,
		}, nil
	}

	// 3. 新文件
	return &response.ResourceCheckResponse{
		Exists: false,
	}, nil
}

// Init 初始化上传任务
func (s *ResourceService) Init(req request.ResourceInit, userUUID uuid.UUID) (*response.ResourceInitResponse, error) {
	// 验证文件类型
	if err := s.getValidator().ValidateMimeType(req.MimeType); err != nil {
		return nil, err
	}

	// 验证文件大小
	if err := s.getValidator().ValidateSize(req.FileSize); err != nil {
		return nil, err
	}

	// 计算总块数
	totalChunks := int(math.Ceil(float64(req.FileSize) / float64(ChunkSize)))
	if totalChunks == 0 {
		totalChunks = 1
	}

	// 初始化QiniuContexts数组
	contexts := make([]string, totalChunks)
	contextsJSON, _ := json.Marshal(contexts)

	// 生成任务ID
	taskUUID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("生成任务ID失败: %w", err)
	}

	// 创建任务
	task := database.ResourceUploadTask{
		TaskID:        taskUUID.String(),
		FileName:      req.FileName,
		FileSize:      req.FileSize,
		FileHash:      req.FileHash,
		MimeType:      req.MimeType,
		ChunkSize:     ChunkSize,
		TotalChunks:   totalChunks,
		Status:        database.TaskStatusInit,
		UserUUID:      userUUID,
		ExpiresAt:     time.Now().Add(TaskExpireHours * time.Hour),
		QiniuContexts: string(contextsJSON),
	}

	if err := global.DB.Create(&task).Error; err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	return &response.ResourceInitResponse{
		TaskID:      task.TaskID,
		TotalChunks: totalChunks,
		ChunkSize:   ChunkSize,
	}, nil
}

// UploadChunkStream 流式上传分片
// reader: 分片数据流
// chunkSize: 分片大小（字节）
func (s *ResourceService) UploadChunkStream(taskID string, chunkNumber int, reader io.Reader, chunkSize int64, userUUID uuid.UUID) (*response.ResourceUploadChunkResponse, error) {
	// 查询任务
	var task database.ResourceUploadTask
	err := global.DB.Where("task_id = ? AND user_uuid = ?", taskID, userUUID).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("查询任务失败: %w", err)
	}

	// 检查任务状态
	if task.Status == database.TaskStatusCompleted {
		return nil, errors.New("任务已完成")
	}
	if task.Status == database.TaskStatusCancelled {
		return nil, errors.New("任务已取消")
	}

	// 检查块号有效性
	if chunkNumber < 0 || chunkNumber >= task.TotalChunks {
		return nil, fmt.Errorf("无效的块号: %d, 有效范围: 0-%d", chunkNumber, task.TotalChunks-1)
	}

	// 流式上传块到七牛云
	ctx, err := s.getUploader().UploadBlockStream(reader, chunkSize)
	if err != nil {
		return nil, fmt.Errorf("上传块失败: %w", err)
	}

	// 使用MySQL JSON_SET原子更新QiniuContexts
	sql := fmt.Sprintf("UPDATE resource_upload_tasks SET qiniu_contexts = JSON_SET(qiniu_contexts, '$[%d]', ?), status = ? WHERE task_id = ?", chunkNumber)
	if err := global.DB.Exec(sql, ctx, database.TaskStatusUploading, taskID).Error; err != nil {
		return nil, fmt.Errorf("更新任务失败: %w", err)
	}

	return &response.ResourceUploadChunkResponse{
		Success:     true,
		ChunkNumber: chunkNumber,
	}, nil
}

// Complete 完成上传
func (s *ResourceService) Complete(req request.ResourceComplete, userUUID uuid.UUID) (*response.ResourceCompleteResponse, error) {
	// 查询任务
	var task database.ResourceUploadTask
	err := global.DB.Where("task_id = ? AND user_uuid = ?", req.TaskID, userUUID).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("查询任务失败: %w", err)
	}

	// 检查任务状态
	if task.Status == database.TaskStatusCompleted {
		return nil, errors.New("任务已完成")
	}

	// 解析contexts
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

	// 生成文件Key
	fileKey := s.getUploader().GenerateFileKey(task.FileName, task.FileHash)

	// 合并文件
	if err := s.getUploader().MergeBlocks(task.FileSize, fileKey, validContexts); err != nil {
		// 更新任务状态为失败
		global.DB.Model(&task).Update("status", database.TaskStatusFailed)
		return nil, fmt.Errorf("合并文件失败: %w", err)
	}

	// 判断是否需要转码（视频文件）
	transcodeStatus := database.TranscodeStatusNone
	if isVideoMimeType(task.MimeType) {
		transcodeStatus = database.TranscodeStatusProcessing
	}

	// 创建资源记录
	resourceRecord := database.Resource{
		FileKey:         fileKey,
		FileName:        task.FileName,
		FileHash:        task.FileHash,
		FileSize:        task.FileSize,
		MimeType:        task.MimeType,
		UserUUID:        userUUID,
		TranscodeStatus: transcodeStatus,
	}
	if err := global.DB.Create(&resourceRecord).Error; err != nil {
		return nil, fmt.Errorf("创建资源记录失败: %w", err)
	}

	// 更新任务状态
	global.DB.Model(&task).Update("status", database.TaskStatusCompleted)

	return &response.ResourceCompleteResponse{
		FileURL: s.getUploader().GetFileURL(fileKey),
		FileKey: fileKey,
	}, nil
}

// Cancel 取消上传
func (s *ResourceService) Cancel(req request.ResourceCancel, userUUID uuid.UUID) error {
	result := global.DB.Model(&database.ResourceUploadTask{}).
		Where("task_id = ? AND user_uuid = ? AND status IN (?, ?)",
			req.TaskID, userUUID, database.TaskStatusInit, database.TaskStatusUploading).
		Update("status", database.TaskStatusCancelled)

	if result.Error != nil {
		return fmt.Errorf("取消任务失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("任务不存在或已完成")
	}
	return nil
}

// Progress 查询上传进度
func (s *ResourceService) Progress(req request.ResourceProgress, userUUID uuid.UUID) (*response.ResourceProgressResponse, error) {
	var task database.ResourceUploadTask
	err := global.DB.Where("task_id = ? AND user_uuid = ?", req.TaskID, userUUID).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("查询任务失败: %w", err)
	}

	uploadedChunks, missingChunks := s.parseContexts(task.QiniuContexts, task.TotalChunks)
	progress := 0
	if task.TotalChunks > 0 {
		progress = len(uploadedChunks) * 100 / task.TotalChunks
	}

	return &response.ResourceProgressResponse{
		TaskID:         task.TaskID,
		TotalChunks:    task.TotalChunks,
		UploadedChunks: uploadedChunks,
		MissingChunks:  missingChunks,
		Progress:       progress,
	}, nil
}

// List 资源列表
func (s *ResourceService) List(req request.ResourceList, userUUID uuid.UUID) ([]response.ResourceItem, int64, error) {
	var resources []database.Resource
	var total int64

	db := global.DB.Model(&database.Resource{}).Where("user_uuid = ?", userUUID)

	if req.FileName != "" {
		db = db.Where("file_name LIKE ?", "%"+req.FileName+"%")
	}
	if req.MimeType != "" {
		db = db.Where("mime_type = ?", req.MimeType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&resources).Error; err != nil {
		return nil, 0, fmt.Errorf("查询列表失败: %w", err)
	}

	uploader := s.getUploader()
	items := make([]response.ResourceItem, len(resources))
	for i, r := range resources {
		item := response.ResourceItem{
			ID:              r.ID,
			FileKey:         r.FileKey,
			FileName:        r.FileName,
			FileURL:         uploader.GetFileURL(r.FileKey),
			FileSize:        r.FileSize,
			MimeType:        r.MimeType,
			TranscodeStatus: int8(r.TranscodeStatus),
		}

		// 如果有转码后的视频，返回转码后的URL
		if r.TranscodeKey != "" {
			item.TranscodeURL = uploader.GetFileURL(r.TranscodeKey)
		}
		// 如果有缩略图，返回缩略图URL
		if r.ThumbnailKey != "" {
			item.ThumbnailURL = uploader.GetFileURL(r.ThumbnailKey)
		}

		items[i] = item
	}

	return items, total, nil
}

// Delete 删除资源
func (s *ResourceService) Delete(req request.ResourceDelete, userUUID uuid.UUID) error {
	// 查询要删除的资源
	var resources []database.Resource
	if err := global.DB.Where("id IN ? AND user_uuid = ?", req.IDs, userUUID).Find(&resources).Error; err != nil {
		return fmt.Errorf("查询资源失败: %w", err)
	}

	if len(resources) == 0 {
		return errors.New("未找到要删除的资源")
	}

	// 删除数据库记录（软删除）
	if err := global.DB.Where("id IN ? AND user_uuid = ?", req.IDs, userUUID).Delete(&database.Resource{}).Error; err != nil {
		return fmt.Errorf("删除资源失败: %w", err)
	}

	// 注意：不删除七牛云上的物理文件，因为可能有其他用户的记录指向同一文件
	// 如果需要删除物理文件，需要检查是否还有其他记录引用该文件

	return nil
}

// parseContexts 解析QiniuContexts，返回已上传和缺失的块号
func (s *ResourceService) parseContexts(contextsJSON string, totalChunks int) ([]int, []int) {
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

// HandleQiniuCallback 处理七牛云转码回调
func (s *ResourceService) HandleQiniuCallback(req request.QiniuCallback) error {
	// 1. 检查整体状态
	if req.Code != 0 {
		// 转码失败，更新状态
		return s.updateTranscodeStatus(req.InputKey, database.TranscodeStatusFailed, "", "")
	}

	// 2. 根据 inputKey 查找数据库记录
	var resource database.Resource
	if err := global.DB.Where("file_key = ?", req.InputKey).First(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("未找到对应的资源记录: %s", req.InputKey)
		}
		return fmt.Errorf("查询资源失败: %w", err)
	}

	// 3. 检查当前状态是否为转码中
	if resource.TranscodeStatus != database.TranscodeStatusProcessing {
		// 非转码中状态，可能是重复回调，忽略
		return nil
	}

	// 4. 遍历 items，根据输出文件后缀提取转码后的 key
	var transcodeKey, thumbnailKey string
	for _, item := range req.Items {
		if item.Code != 0 {
			continue // 跳过失败的项
		}

		key := item.Key
		// 根据输出文件后缀判断类型
		if strings.HasSuffix(key, "_h264.mp4") {
			transcodeKey = key
		} else if strings.HasSuffix(key, "_thumb.jpg") {
			thumbnailKey = key
		}
	}

	// 5. 更新数据库
	return s.updateTranscodeStatus(req.InputKey, database.TranscodeStatusSuccess, transcodeKey, thumbnailKey)
}

// updateTranscodeStatus 更新转码状态
func (s *ResourceService) updateTranscodeStatus(fileKey string, status database.TranscodeStatus, transcodeKey, thumbnailKey string) error {
	updates := map[string]interface{}{
		"transcode_status": status,
	}

	if transcodeKey != "" {
		updates["transcode_key"] = transcodeKey
	}
	if thumbnailKey != "" {
		updates["thumbnail_key"] = thumbnailKey
	}

	result := global.DB.Model(&database.Resource{}).Where("file_key = ?", fileKey).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新转码状态失败: %w", result.Error)
	}

	return nil
}

// isVideoMimeType 判断是否为视频 MIME 类型
func isVideoMimeType(mimeType string) bool {
	return len(mimeType) >= 6 && mimeType[:6] == "video/"
}
