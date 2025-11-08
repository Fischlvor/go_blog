package api

import (
	"server/internal/model/request"
	"server/internal/model/response"
	"server/pkg/global"
	"server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type AIChatApi struct{}

// CreateSession 创建聊天会话
func (api *AIChatApi) CreateSession(c *gin.Context) {
	var req request.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// ✅ 获取当前用户UUID
	userUUID := utils.GetUUID(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户未登录", c)
		return
	}

	session, err := aiChatService.CreateSession(userUUID, req)
	if err != nil {
		global.Log.Error("创建会话失败: " + err.Error())
		response.FailWithMessage("创建会话失败", c)
		return
	}

	response.OkWithData(session, c)
}

// GetSessions 获取会话列表
func (api *AIChatApi) GetSessions(c *gin.Context) {
	var req request.GetSessionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// ✅ 获取当前用户UUID
	userUUID := utils.GetUUID(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户未登录", c)
		return
	}

	sessions, total, err := aiChatService.GetSessions(userUUID, req)
	if err != nil {
		global.Log.Error("获取会话列表失败: " + err.Error())
		response.FailWithMessage("获取会话列表失败", c)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  sessions,
		"total": total,
	}, "获取成功", c)
}

// GetMessages 获取消息列表
func (api *AIChatApi) GetMessages(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		response.FailWithMessage("会话ID不能为空", c)
		return
	}

	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("会话ID格式错误", c)
		return
	}

	var req request.GetMessagesRequest
	req.SessionID = uint(sessionID)
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 获取当前用户ID
	userUUID := utils.GetUUID(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户未登录", c)
		return
	}

	messages, total, err := aiChatService.GetMessages(userUUID, req)
	if err != nil {
		global.Log.Error("获取消息列表失败: " + err.Error())
		response.FailWithMessage("获取消息列表失败", c)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  messages,
		"total": total,
	}, "获取成功", c)
}

// SendMessage 发送消息
func (api *AIChatApi) SendMessage(c *gin.Context) {
	var req request.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 获取当前用户ID
	userUUID := utils.GetUUID(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户未登录", c)
		return
	}

	chatResp, err := aiChatService.SendMessage(userUUID, req)
	if err != nil {
		global.Log.Error("发送消息失败: " + err.Error())
		response.FailWithMessage("发送消息失败", c)
		return
	}

	response.OkWithData(chatResp, c)
}

// SendMessageStream 流式发送消息
func (api *AIChatApi) SendMessageStream(c *gin.Context) {
	var req request.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 获取当前用户ID
	userUUID := utils.GetUUID(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户未登录", c)
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 创建一个通道用于通知客户端连接关闭
	notify := c.Writer.CloseNotify()

	// 创建一个自定义的writer来处理SSE格式
	sseWriter := &SSEWriter{writer: c.Writer}

	err := aiChatService.SendMessageStream(userUUID, req, sseWriter)
	if err != nil {
		global.Log.Error("流式发送消息失败: " + err.Error())
		// 发送错误消息
		c.SSEvent("error", gin.H{"message": "发送消息失败"})
		return
	}

	// 监听客户端断开连接
	<-notify
}

// DeleteSession 删除会话
func (api *AIChatApi) DeleteSession(c *gin.Context) {
	var req request.DeleteSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 获取当前用户ID
	userUUID := utils.GetUUID(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户未登录", c)
		return
	}

	err := aiChatService.DeleteSession(userUUID, req)
	if err != nil {
		global.Log.Error("删除会话失败: " + err.Error())
		response.FailWithMessage("删除会话失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// UpdateSession 更新会话
func (api *AIChatApi) UpdateSession(c *gin.Context) {
	var req request.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 获取当前用户ID
	userUUID := utils.GetUUID(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户未登录", c)
		return
	}

	err := aiChatService.UpdateSession(userUUID, req)
	if err != nil {
		global.Log.Error("更新会话失败: " + err.Error())
		response.FailWithMessage("更新会话失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// GetSessionDetail 获取会话详情
func (api *AIChatApi) GetSessionDetail(c *gin.Context) {
	sessionIDStr := c.Param("id")
	if sessionIDStr == "" {
		response.FailWithMessage("会话ID不能为空", c)
		return
	}

	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("会话ID格式错误", c)
		return
	}

	// 获取当前用户ID
	userUUID := utils.GetUUID(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户未登录", c)
		return
	}

	session, err := aiChatService.GetSessionDetail(userUUID, uint(sessionID))
	if err != nil {
		global.Log.Error("获取会话详情失败: " + err.Error())
		response.FailWithMessage("获取会话详情失败", c)
		return
	}

	response.OkWithData(session, c)
}

// GetAvailableModels 获取可用模型列表
func (api *AIChatApi) GetAvailableModels(c *gin.Context) {
	models, err := aiChatService.GetAvailableModels()
	if err != nil {
		global.Log.Error("获取模型列表失败: " + err.Error())
		response.FailWithMessage("获取模型列表失败", c)
		return
	}

	response.OkWithData(models, c)
}

// SSEWriter SSE响应写入器
type SSEWriter struct {
	writer gin.ResponseWriter
}

func (w *SSEWriter) Write(p []byte) (n int, err error) {
	// 直接写入数据，因为服务层已经包装了SSE格式
	return w.writer.Write(p)
}

func (w *SSEWriter) Flush() {
	if flusher, ok := w.writer.(interface{ Flush() }); ok {
		flusher.Flush()
	}
}
