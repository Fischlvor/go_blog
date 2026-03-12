import { adminService } from '@/utils/request'

// AI模型管理
export const aiModelApi = {
  // 创建AI模型
  create: (data: any) => adminService.post('/ai-management/model', data),
  
  // 删除AI模型
  delete: (id: number) => adminService.delete(`/ai-management/model/${id}`),
  
  // 更新AI模型
  update: (data: any) => adminService.put('/ai-management/model', data),
  
  // 获取AI模型详情
  get: (id: number) => adminService.get(`/ai-management/model/${id}`),
  
  // 获取AI模型列表
  list: (params: any) => adminService.get('/ai-management/model', { params })
}

// AI会话管理
export const aiSessionApi = {
  // 获取AI会话列表
  list: (params: any) => adminService.get('/ai-management/session', { params }),
  
  // 获取AI会话详情
  get: (id: number) => adminService.get(`/ai-management/session/${id}`),
  
  // 删除AI会话
  delete: (id: number) => adminService.delete(`/ai-management/session/${id}`)
}

// AI消息管理
export const aiMessageApi = {
  // 获取AI消息列表
  list: (params: any) => adminService.get('/ai-management/message', { params }),
  
  // 获取AI消息详情
  get: (id: number) => adminService.get(`/ai-management/message/${id}`),
  
  // 删除AI消息
  delete: (id: number) => adminService.delete(`/ai-management/message/${id}`)
} 