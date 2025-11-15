package admin

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"server/internal/model"
	"server/internal/model/response"
	"server/internal/service"
	"server/pkg/global"
	"server/pkg/utils"
	"server/pkg/utils/upload"
)

type EmojiApi struct {
	emojiService *service.EmojiService
}

func NewEmojiApi() *EmojiApi {
	return &EmojiApi{
		emojiService: &service.EmojiService{},
	}
}

// GetEmojiList 获取表情列表
// @Summary 获取表情列表
// @Description 分页获取表情列表，支持按组筛选和关键词搜索
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param group_key query string false "组key筛选"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} response.Response{data=model.EmojiListResponse}
// @Router /api/admin/emoji/list [get]
func (a *EmojiApi) GetEmojiList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	groupKey := c.Query("group_key")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := a.emojiService.GetEmojiList(page, pageSize, groupKey, keyword)
	if err != nil {
		response.FailWithMessage("获取表情列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetEmojiGroups 获取表情组列表
// @Summary 获取表情组列表
// @Description 获取所有表情组
// @Tags 表情管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]model.EmojiGroup}
// @Router /api/admin/emoji/groups [get]
func (a *EmojiApi) GetEmojiGroups(c *gin.Context) {
	groups, err := a.emojiService.GetEmojiGroups()
	if err != nil {
		response.FailWithMessage("获取表情组失败: "+err.Error(), c)
		return
	}

	response.OkWithData(groups, c)
}

// CreateEmojiGroup 创建表情组
// @Summary 创建表情组
// @Description 创建新的表情组
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param body body object{group_name=string,description=string} true "表情组信息"
// @Success 200 {object} response.Response
// @Router /api/admin/emoji/groups [post]
func (a *EmojiApi) CreateEmojiGroup(c *gin.Context) {
	var req struct {
		GroupName   string `json:"group_name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 获取当前用户UUID
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		response.NoAuth("未登录", c)
		return
	}

	err := a.emojiService.CreateEmojiGroup(req.GroupName, req.Description, userUUID.(uuid.UUID))
	if err != nil {
		response.FailWithMessage("创建表情组失败: "+err.Error(), c)
		return
	}

	response.OkWithData("创建成功", c)
}

// UpdateEmojiGroup 更新表情组
// @Summary 更新表情组
// @Description 更新表情组信息
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param id path int true "组ID"
// @Param body body object{group_name=string,description=string} true "表情组信息"
// @Success 200 {object} response.Response
// @Router /api/admin/emoji/groups/{id} [put]
func (a *EmojiApi) UpdateEmojiGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("无效的组ID", c)
		return
	}

	var req struct {
		GroupName   string `json:"group_name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	err = a.emojiService.UpdateEmojiGroup(id, req.GroupName, req.Description)
	if err != nil {
		response.FailWithMessage("更新表情组失败: "+err.Error(), c)
		return
	}

	response.OkWithData("更新成功", c)
}

// DeleteEmojiGroup 删除表情组
// @Summary 删除表情组
// @Description 删除表情组（需要确保组内无表情）
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param id path int true "组ID"
// @Success 200 {object} response.Response
// @Router /api/admin/emoji/groups/{id} [delete]
func (a *EmojiApi) DeleteEmojiGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("无效的组ID", c)
		return
	}

	err = a.emojiService.DeleteEmojiGroup(id)
	if err != nil {
		response.FailWithMessage("删除表情组失败: "+err.Error(), c)
		return
	}

	response.OkWithData("删除成功", c)
}

// UploadEmoji 上传表情
// @Summary 批量上传表情
// @Description 批量上传表情文件
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param body body model.EmojiUploadRequest true "上传请求"
// @Success 200 {object} response.Response{data=[]model.Emoji}
// @Router /api/admin/emoji/upload [post]
func (a *EmojiApi) UploadEmoji(c *gin.Context) {
	// 获取当前用户UUID
	userUUID := utils.GetUUID(c)
	if userUUID == (uuid.UUID{}) {
		response.NoAuth("未登录", c)
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 流式处理multipart上传
	a.handleMultipartUploadStream(c, userUUID)
}

// DeleteEmoji 删除表情
// @Summary 删除表情
// @Description 软删除表情
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param id path int true "表情ID"
// @Success 200 {object} response.Response
// @Router /api/admin/emoji/{id} [delete]
func (a *EmojiApi) DeleteEmoji(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("无效的表情ID", c)
		return
	}

	err = a.emojiService.DeleteEmoji(id)
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}

	response.OkWithData("删除成功", c)
}

// RestoreEmoji 恢复表情
// @Summary 恢复表情
// @Description 恢复已删除的表情
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param id path int true "表情ID"
// @Success 200 {object} response.Response
// @Router /api/admin/emoji/{id}/restore [put]
func (a *EmojiApi) RestoreEmoji(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("无效的表情ID", c)
		return
	}

	err = a.emojiService.RestoreEmoji(id)
	if err != nil {
		response.FailWithMessage("恢复失败: "+err.Error(), c)
		return
	}

	response.OkWithData("恢复成功", c)
}

// RegenerateSprites 重新生成雪碧图（SSE流式响应）
// @Summary 重新生成雪碧图
// @Description 重新生成雪碧图，支持指定表情组或全部生成
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param body body object{group_keys=[]string} false "表情组keys数组，为空则生成全部"
// @Success 200 {object} response.Response
// @Router /api/admin/emoji/regenerate [post]
func (a *EmojiApi) RegenerateSprites(c *gin.Context) {
	// 获取当前用户UUID
	userUUID := utils.GetUUID(c)
	if userUUID == (uuid.UUID{}) {
		response.NoAuth("未登录", c)
		return
	}

	// 解析请求体
	var req struct {
		GroupKeys []string `json:"group_keys"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果没有body，则默认为空数组（生成全部）
		req.GroupKeys = []string{}
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 创建SSE写入器
	sseWriter := &APISSEWriter{writer: c.Writer}

	// 直接执行雪碧图生成（不使用异步任务）
	err := a.emojiService.DoRegenerateSpritesWithSSE(userUUID, req.GroupKeys, sseWriter)
	if err != nil {
		sseWriter.WriteSSEEvent("error", gin.H{
			"message": fmt.Sprintf("雪碧图生成失败: %v", err),
		})
		return
	}
}

// APISSEWriter API层的SSE写入器实现
type APISSEWriter struct {
	writer gin.ResponseWriter
}

func (w *APISSEWriter) WriteSSEEvent(event string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.writer.WriteString(fmt.Sprintf("event: %s\n", event))
	w.writer.WriteString(fmt.Sprintf("data: %s\n\n", string(jsonData)))
	w.Flush()
	return nil
}

func (w *APISSEWriter) Flush() {
	if flusher, ok := w.writer.(interface{ Flush() }); ok {
		flusher.Flush()
	}
}

// GetTaskStatus 获取任务状态
// @Summary 获取任务状态
// @Description 获取异步任务的执行状态
// @Tags 表情管理
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response{data=model.EmojiTask}
// @Router /api/admin/emoji/task/{id} [get]
func (a *EmojiApi) GetTaskStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("无效的任务ID", c)
		return
	}

	task, err := a.emojiService.GetTaskStatus(id)
	if err != nil {
		response.FailWithMessage("获取任务状态失败: "+err.Error(), c)
		return
	}

	response.OkWithData(task, c)
}

// GetEmojiConfig 获取前端配置
// @Summary 获取前端配置
// @Description 获取前端使用的emoji配置信息
// @Tags 表情管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=model.EmojiConfigResponse}
// @Router /api/emoji/config [get]
func (a *EmojiApi) GetEmojiConfig(c *gin.Context) {
	config, err := a.emojiService.GetEmojiConfig()
	if err != nil {
		response.FailWithMessage("获取配置失败: "+err.Error(), c)
		return
	}

	response.OkWithData(config, c)
}

// GetSpriteList 获取雪碧图列表
// @Summary 获取雪碧图列表
// @Description 获取所有雪碧图的列表信息
// @Tags 表情管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]model.EmojiSprite}
// @Router /api/admin/emoji/sprites [get]
func (a *EmojiApi) GetSpriteList(c *gin.Context) {
	sprites, err := a.emojiService.GetSpriteList()
	if err != nil {
		response.FailWithMessage("获取雪碧图列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(sprites, c)
}

// handleMultipartUpload 处理multipart格式的上传
func (a *EmojiApi) handleMultipartUpload(c *gin.Context, userUUID uuid.UUID) ([]model.Emoji, error) {
	// 获取group_key参数
	groupKey := c.PostForm("group_key")
	if groupKey == "" {
		return nil, fmt.Errorf("group_key参数不能为空")
	}

	// 获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("解析multipart表单失败: %v", err)
	}

	files := form.File["files"]
	if len(files) == 0 {
		return nil, fmt.Errorf("没有上传的文件")
	}

	// 使用现有的上传系统处理文件
	var emojis []model.Emoji
	oss := upload.NewOss()

	// 查询表情组信息
	var group model.EmojiGroup
	if err := global.DB.Where("group_key = ?", groupKey).First(&group).Error; err != nil {
		return nil, fmt.Errorf("表情组不存在: %s", groupKey)
	}

	for i, fileHeader := range files {
		// 使用带路径前缀的上传方法
		pathPrefix := fmt.Sprintf("emoji/%s", groupKey)
		url, filename, err := oss.UploadImageWithPrefix(fileHeader, pathPrefix)
		if err != nil {
			return nil, fmt.Errorf("文件 %s 上传失败: %v", fileHeader.Filename, err)
		}

		// 生成key和计算雪碧图位置
		emojiKey := fmt.Sprintf("e%d", group.EmojiCount+i)
		keyNum := group.EmojiCount + i
		spriteGroup := keyNum / 128
		posInSprite := keyNum % 128
		row := posInSprite / 16
		col := posInSprite % 16

		// 创建emoji记录并保存到数据库
		emoji := model.Emoji{
			Key:             emojiKey,
			Filename:        filename,
			GroupKey:        group.GroupKey,
			SpriteGroup:     spriteGroup,
			SpritePositionX: col * 32,
			SpritePositionY: row * 32,
			CdnUrl:          url,
			FileSize:        int(fileHeader.Size),
			Status:          model.EmojiStatusActive,
			CreatedBy:       userUUID,
		}

		// 保存到数据库
		if err := global.DB.Create(&emoji).Error; err != nil {
			return nil, fmt.Errorf("保存表情记录失败: %v", err)
		}

		emojis = append(emojis, emoji)
	}

	// 更新表情组的emoji_count
	if err := global.DB.Model(&group).Update("emoji_count", group.EmojiCount+len(files)).Error; err != nil {
		return nil, fmt.Errorf("更新表情组计数失败: %v", err)
	}

	return emojis, nil
}

// handleMultipartUploadStream 流式处理multipart格式的上传
func (a *EmojiApi) handleMultipartUploadStream(c *gin.Context, userUUID uuid.UUID) {
	// 获取group_key参数
	groupKey := c.PostForm("group_key")
	if groupKey == "" {
		a.sendSSEError(c, "group_key参数不能为空")
		return
	}

	// 获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		a.sendSSEError(c, fmt.Sprintf("解析multipart表单失败: %v", err))
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		a.sendSSEError(c, "没有上传的文件")
		return
	}

	// 查询表情组信息
	var group model.EmojiGroup
	if err := global.DB.Where("group_key = ?", groupKey).First(&group).Error; err != nil {
		a.sendSSEError(c, fmt.Sprintf("表情组不存在: %s", groupKey))
		return
	}

	// 发送开始事件
	a.sendSSEEvent(c, "start", gin.H{
		"total":   len(files),
		"message": "开始处理文件上传...",
	})

	// 使用现有的上传系统处理文件
	var emojis []model.Emoji
	var failedFiles []string
	oss := upload.NewOss()

	for i, fileHeader := range files {
		// 发送当前处理进度
		a.sendSSEEvent(c, "progress", gin.H{
			"current":  i + 1,
			"total":    len(files),
			"filename": fileHeader.Filename,
			"message":  fmt.Sprintf("正在处理文件: %s", fileHeader.Filename),
		})

		// 使用带路径前缀的上传方法
		pathPrefix := fmt.Sprintf("emoji/%s", groupKey)
		url, filename, err := oss.UploadImageWithPrefix(fileHeader, pathPrefix)
		if err != nil {
			a.sendSSEEvent(c, "file_error", gin.H{
				"filename": fileHeader.Filename,
				"error":    fmt.Sprintf("文件上传失败: %v", err),
			})
			failedFiles = append(failedFiles, fileHeader.Filename)
			continue
		}

		// 生成key和计算雪碧图位置
		emojiKey := fmt.Sprintf("e%d", group.EmojiCount+i)
		keyNum := group.EmojiCount + i
		spriteGroup := keyNum / 128
		posInSprite := keyNum % 128
		row := posInSprite / 16
		col := posInSprite % 16

		// 创建emoji记录（暂不保存到数据库）
		emoji := model.Emoji{
			Key:             emojiKey,
			Filename:        filename,
			GroupKey:        group.GroupKey,
			SpriteGroup:     spriteGroup,
			SpritePositionX: col * 32,
			SpritePositionY: row * 32,
			CdnUrl:          url,
			FileSize:        int(fileHeader.Size),
			Status:          model.EmojiStatusActive,
			CreatedBy:       userUUID,
		}

		emojis = append(emojis, emoji)

		// 发送单个文件完成事件
		a.sendSSEEvent(c, "file_complete", gin.H{
			"filename": fileHeader.Filename,
			"current":  i + 1,
			"total":    len(files),
		})
	}

	// 批量保存所有emoji到数据库
	if len(emojis) > 0 {
		if err := global.DB.CreateInBatches(&emojis, 100).Error; err != nil {
			a.sendSSEEvent(c, "error", gin.H{
				"message": fmt.Sprintf("批量保存表情记录失败: %v", err),
			})
			return
		}
	}

	// 更新表情组的emoji_count
	if err := global.DB.Model(&group).Update("emoji_count", group.EmojiCount+len(emojis)).Error; err != nil {
		a.sendSSEEvent(c, "warning", gin.H{
			"message": fmt.Sprintf("更新表情组计数失败: %v", err),
		})
	}

	// 发送完成事件
	a.sendSSEEvent(c, "complete", gin.H{
		"message": "所有文件处理完成",
		"total":   len(files),
		"success": len(emojis),
		"failed":  len(failedFiles),
	})
}

// sendSSEEvent 发送SSE事件
func (a *EmojiApi) sendSSEEvent(c *gin.Context, event string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	c.Writer.WriteString(fmt.Sprintf("event: %s\n", event))
	c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", string(jsonData)))
	c.Writer.Flush()
}

// sendSSEError 发送SSE错误事件
func (a *EmojiApi) sendSSEError(c *gin.Context, message string) {
	a.sendSSEEvent(c, "error", gin.H{"message": message})
}
