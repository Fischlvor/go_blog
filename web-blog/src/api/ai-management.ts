import request from '@/utils/request'

// AI模型管理
export const aiModelApi = {
  // 创建AI模型
  create: (data: any) => request.post('/ai-management/model', data),
  
  // 删除AI模型
  delete: (id: number) => request.delete(`/ai-management/model/${id}`),
  
  // 更新AI模型
  update: (data: any) => request.put('/ai-management/model', data),
  
  // 获取AI模型详情
  get: (id: number) => request.get(`/ai-management/model/${id}`),
  
  // 获取AI模型列表
  list: (params: any) => request.get('/ai-management/model', { params })
}

// AI会话管理
export const aiSessionApi = {
  // 获取AI会话列表
  list: (params: any) => request.get('/ai-management/session', { params }),
  
  // 获取AI会话详情
  get: (id: number) => request.get(`/ai-management/session/${id}`),
  
  // 删除AI会话
  delete: (id: number) => request.delete(`/ai-management/session/${id}`)
}

// AI消息管理
export const aiMessageApi = {
  // 获取AI消息列表
  list: (params: any) => request.get('/ai-management/message', { params }),
  
  // 获取AI消息详情
  get: (id: number) => request.get(`/ai-management/message/${id}`),
  
  // 删除AI消息
  delete: (id: number) => request.delete(`/ai-management/message/${id}`)
} 