package api

import (
	"server/pkg/global"
	"server/internal/model/database"
	"server/internal/model/request"
	"server/internal/model/response"
	"server/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AIManagementApi struct{}

// ==================== AI模型管理 ====================

// CreateAIModel 创建AI模型
func (a *AIManagementApi) CreateAIModel(c *gin.Context) {
	var model database.AIModel
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err = service.ServiceGroupApp.AIManagementService.CreateAIModel(model)
	if err != nil {
		global.Log.Error("创建AI模型失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteAIModel 删除AI模型
func (a *AIManagementApi) DeleteAIModel(c *gin.Context) {
	id := c.Param("id")
	modelID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err = service.ServiceGroupApp.AIManagementService.DeleteAIModel(uint(modelID))
	if err != nil {
		global.Log.Error("删除AI模型失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateAIModel 更新AI模型
func (a *AIManagementApi) UpdateAIModel(c *gin.Context) {
	var model database.AIModel
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err = service.ServiceGroupApp.AIManagementService.UpdateAIModel(model)
	if err != nil {
		global.Log.Error("更新AI模型失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetAIModel 获取AI模型详情
func (a *AIManagementApi) GetAIModel(c *gin.Context) {
	id := c.Param("id")
	modelID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	model, err := service.ServiceGroupApp.AIManagementService.GetAIModel(uint(modelID))
	if err != nil {
		global.Log.Error("获取AI模型失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(model, c)
}

// GetAIModelList 获取AI模型列表
func (a *AIManagementApi) GetAIModelList(c *gin.Context) {
	var req request.AIModelListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	list, total, err := service.ServiceGroupApp.AIManagementService.GetAIModelList(req)
	if err != nil {
		global.Log.Error("获取AI模型列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "获取成功", c)
}

// ==================== AI会话管理 ====================

// GetAISessionList 获取AI会话列表
func (a *AIManagementApi) GetAISessionList(c *gin.Context) {
	var req request.AISessionListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	list, total, err := service.ServiceGroupApp.AIManagementService.GetAISessionList(req)
	if err != nil {
		global.Log.Error("获取AI会话列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "获取成功", c)
}

// GetAISession 获取AI会话详情
func (a *AIManagementApi) GetAISession(c *gin.Context) {
	id := c.Param("id")
	sessionID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	session, err := service.ServiceGroupApp.AIManagementService.GetAISession(uint(sessionID))
	if err != nil {
		global.Log.Error("获取AI会话失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(session, c)
}

// DeleteAISession 删除AI会话
func (a *AIManagementApi) DeleteAISession(c *gin.Context) {
	id := c.Param("id")
	sessionID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err = service.ServiceGroupApp.AIManagementService.DeleteAISession(uint(sessionID))
	if err != nil {
		global.Log.Error("删除AI会话失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ==================== AI消息管理 ====================

// GetAIMessageList 获取AI消息列表
func (a *AIManagementApi) GetAIMessageList(c *gin.Context) {
	var req request.AIMessageListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	list, total, err := service.ServiceGroupApp.AIManagementService.GetAIMessageList(req)
	if err != nil {
		global.Log.Error("获取AI消息列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "获取成功", c)
}

// GetAIMessage 获取AI消息详情
func (a *AIManagementApi) GetAIMessage(c *gin.Context) {
	id := c.Param("id")
	messageID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	message, err := service.ServiceGroupApp.AIManagementService.GetAIMessage(uint(messageID))
	if err != nil {
		global.Log.Error("获取AI消息失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(message, c)
}

// DeleteAIMessage 删除AI消息
func (a *AIManagementApi) DeleteAIMessage(c *gin.Context) {
	id := c.Param("id")
	messageID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err = service.ServiceGroupApp.AIManagementService.DeleteAIMessage(uint(messageID))
	if err != nil {
		global.Log.Error("删除AI消息失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
