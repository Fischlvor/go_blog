package api

import (
	"server/internal/model/request"
	"server/internal/model/response"
	"server/internal/service"
	"server/pkg/global"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ResourceApi struct{}

var resourceService = &service.ResourceService{}

// Check 检查文件（秒传/续传检测）
// POST /api/admin/resources/check
func (r *ResourceApi) Check(c *gin.Context) {
	var req request.ResourceCheck
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userUUID := utils.GetUUID(c)
	result, err := resourceService.Check(req, userUUID)
	if err != nil {
		global.Log.Error("检查文件失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// Init 初始化上传任务
// POST /api/admin/resources/init
func (r *ResourceApi) Init(c *gin.Context) {
	var req request.ResourceInit
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userUUID := utils.GetUUID(c)
	result, err := resourceService.Init(req, userUUID)
	if err != nil {
		global.Log.Error("初始化任务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// UploadChunk 上传分片（流式）
// POST /api/admin/resources/upload-chunk
func (r *ResourceApi) UploadChunk(c *gin.Context) {
	// 获取任务ID
	taskID := c.PostForm("task_id")
	if taskID == "" {
		response.FailWithMessage("task_id不能为空", c)
		return
	}

	// 获取块号
	chunkNumberStr := c.PostForm("chunk_number")
	var chunkNumber int
	if _, err := parseChunkNumber(chunkNumberStr, &chunkNumber); err != nil {
		response.FailWithMessage("chunk_number无效", c)
		return
	}

	// 获取块数据（流式，不读入内存）
	file, header, err := c.Request.FormFile("chunk_data")
	if err != nil {
		response.FailWithMessage("获取块数据失败: "+err.Error(), c)
		return
	}
	defer file.Close()

	// 获取块大小
	chunkSize := header.Size

	userUUID := utils.GetUUID(c)
	// 流式上传：直接传递 io.Reader，不读入内存
	result, err := resourceService.UploadChunkStream(taskID, chunkNumber, file, chunkSize, userUUID)
	if err != nil {
		global.Log.Error("上传分片失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// Complete 完成上传
// POST /api/admin/resources/complete
func (r *ResourceApi) Complete(c *gin.Context) {
	var req request.ResourceComplete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userUUID := utils.GetUUID(c)
	result, err := resourceService.Complete(req, userUUID)
	if err != nil {
		global.Log.Error("完成上传失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// Cancel 取消上传
// POST /api/admin/resources/cancel
func (r *ResourceApi) Cancel(c *gin.Context) {
	var req request.ResourceCancel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userUUID := utils.GetUUID(c)
	err := resourceService.Cancel(req, userUUID)
	if err != nil {
		global.Log.Error("取消上传失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("取消成功", c)
}

// Progress 查询上传进度
// GET /api/admin/resources/progress
func (r *ResourceApi) Progress(c *gin.Context) {
	var req request.ResourceProgress
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userUUID := utils.GetUUID(c)
	result, err := resourceService.Progress(req, userUUID)
	if err != nil {
		global.Log.Error("查询进度失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// List 资源列表
// GET /api/admin/resources/list
func (r *ResourceApi) List(c *gin.Context) {
	var req request.ResourceList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userUUID := utils.GetUUID(c)
	items, total, err := resourceService.List(req, userUUID)
	if err != nil {
		global.Log.Error("查询资源列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  items,
		Total: total,
	}, c)
}

// Delete 删除资源
// POST /api/admin/resources/delete
func (r *ResourceApi) Delete(c *gin.Context) {
	var req request.ResourceDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userUUID := utils.GetUUID(c)
	err := resourceService.Delete(req, userUUID)
	if err != nil {
		global.Log.Error("删除资源失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// parseChunkNumber 解析块号
func parseChunkNumber(s string, result *int) (bool, error) {
	if s == "" {
		return false, nil
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, nil
		}
		n = n*10 + int(c-'0')
	}
	*result = n
	return true, nil
}

// QiniuCallback 七牛云转码回调
// POST /api/qiniu/callback
// 此接口不需要认证，由七牛云服务器调用
func (r *ResourceApi) QiniuCallback(c *gin.Context) {
	var req request.QiniuCallback
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("解析七牛回调请求失败", zap.Error(err))
		// 返回成功，避免七牛重试（解析失败重试也没用）
		c.JSON(200, gin.H{"code": 0, "msg": "ok"})
		return
	}

	// 打印回调信息（用于调试）
	global.Log.Info("七牛云转码回调",
		zap.String("id", req.ID),
		zap.Int("code", req.Code),
		zap.String("inputKey", req.InputKey),
		zap.Int("items_count", len(req.Items)),
	)

	// 处理回调
	if err := resourceService.HandleQiniuCallback(req); err != nil {
		global.Log.Error("处理七牛回调失败",
			zap.Error(err),
			zap.String("inputKey", req.InputKey),
		)
		// 仍然返回成功，避免无限重试
		// 错误已记录日志，可以后续手动处理
	}

	// 返回成功，否则七牛会重试
	c.JSON(200, gin.H{"code": 0, "msg": "ok"})
}
