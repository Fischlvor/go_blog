import request, { streamRequest, type StreamCallbacks } from '@/utils/request'

// AI聊天相关接口
export const aiChatApi = {
  // 获取可用模型列表
  getAvailableModels() {
    return request({
      url: '/ai-chat/models',
      method: 'get'
    })
  },

  // 创建会话
  createSession(data: { title: string; model: string }) {
    return request({
      url: '/ai-chat/session',
      method: 'post',
      data
    })
  },

  // 获取会话列表
  getSessions(params?: { page?: number; pageSize?: number }) {
    return request({
      url: '/ai-chat/sessions',
      method: 'get',
      params
    })
  },

  // 获取消息列表
  getMessages(sessionId: number, params?: { page?: number; pageSize?: number }) {
    return request({
      url: '/ai-chat/messages',
      method: 'get',
      params: { session_id: sessionId, ...params }
    })
  },

  // 发送消息 - 设置20秒超时时间
  sendMessage(data: { session_id: number; content: string }) {
    return request({
      url: '/ai-chat/message',
      method: 'post',
      data,
      timeout: 20000 // 20秒超时
    })
  },

  // 流式发送消息 - 使用通用流式请求方法
  sendMessageStream(data: { session_id: number; content: string }, callbacks: StreamCallbacks) {
    return streamRequest({
      url: '/ai-chat/message/stream',
      method: 'POST',
      data,
      timeout: 60000 // 60秒超时
    }, callbacks)
  },

  // 删除会话
  deleteSession(data: { session_id: number }) {
    return request({
      url: '/ai-chat/session',
      method: 'delete',
      data
    })
  },

  // 更新会话
  updateSession(data: { session_id: number; title: string }) {
    return request({
      url: '/ai-chat/session',
      method: 'put',
      data
    })
  }
} 