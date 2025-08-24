package service

import (
	"server/global"
	"server/model/database"
	"server/model/request"

	"gorm.io/gorm"
)

type AIManagementService struct{}

// ==================== AI模型管理 ====================

// CreateAIModel 创建AI模型
func (s *AIManagementService) CreateAIModel(model database.AIModel) error {
	return global.DB.Create(&model).Error
}

// DeleteAIModel 删除AI模型
func (s *AIManagementService) DeleteAIModel(id uint) error {
	return global.DB.Delete(&database.AIModel{}, id).Error
}

// UpdateAIModel 更新AI模型
func (s *AIManagementService) UpdateAIModel(model database.AIModel) error {
	return global.DB.Save(&model).Error
}

// GetAIModel 获取AI模型详情
func (s *AIManagementService) GetAIModel(id uint) (database.AIModel, error) {
	var model database.AIModel
	err := global.DB.First(&model, id).Error
	return model, err
}

// GetAIModelList 获取AI模型列表
func (s *AIManagementService) GetAIModelList(info request.AIModelListRequest) (list []database.AIModel, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&database.AIModel{})

	// 添加查询条件
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Provider != "" {
		db = db.Where("provider LIKE ?", "%"+info.Provider+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return
}

// ==================== AI会话管理 ====================

// GetAISessionList 获取AI会话列表
func (s *AIManagementService) GetAISessionList(info request.AISessionListRequest) (list []database.AIChatSession, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&database.AIChatSession{})

	// 添加查询条件
	if info.UserID != "" {
		db = db.Where("user_id = ?", info.UserID)
	}
	if info.Model != "" {
		db = db.Where("model = ?", info.Model)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Preload("User").Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return
}

// GetAISession 获取AI会话详情
func (s *AIManagementService) GetAISession(id uint) (database.AIChatSession, error) {
	var session database.AIChatSession
	err := global.DB.Preload("User").Preload("Messages").First(&session, id).Error
	return session, err
}

// DeleteAISession 删除AI会话
func (s *AIManagementService) DeleteAISession(id uint) error {
	// 开启事务
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 先删除会话下的所有消息
		if err := tx.Where("session_id = ?", id).Delete(&database.AIChatMessage{}).Error; err != nil {
			return err
		}
		// 再删除会话
		return tx.Delete(&database.AIChatSession{}, id).Error
	})
}

// ==================== AI消息管理 ====================

// GetAIMessageList 获取AI消息列表
func (s *AIManagementService) GetAIMessageList(info request.AIMessageListRequest) (list []database.AIChatMessage, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&database.AIChatMessage{})

	// 添加查询条件
	if info.SessionID != "" {
		db = db.Where("session_id = ?", info.SessionID)
	}
	if info.Role != "" {
		db = db.Where("role = ?", info.Role)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return
}

// GetAIMessage 获取AI消息详情
func (s *AIManagementService) GetAIMessage(id uint) (database.AIChatMessage, error) {
	var message database.AIChatMessage
	err := global.DB.First(&message, id).Error
	return message, err
}

// DeleteAIMessage 删除AI消息
func (s *AIManagementService) DeleteAIMessage(id uint) error {
	return global.DB.Delete(&database.AIChatMessage{}, id).Error
}
