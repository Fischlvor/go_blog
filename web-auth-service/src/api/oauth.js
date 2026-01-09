import request from '../utils/request'

// 获取公开应用列表
export const getPublicApplications = () => {
  return request.get('/oauth/applications')
}
