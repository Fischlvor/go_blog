import request from '../utils/request'

export const manageApi = {
  // 获取设备列表
  getDevices: () => request.get('/manage/devices'),
  
  // 踢出设备
  kickDevice: (deviceId) => request.post('/manage/kick-device', { device_id: deviceId }),
  
  // 退出manage应用
  logout: () => request.post('/manage/logout'),
  
  // SSO全局退出
  ssoLogout: () => request.post('/manage/sso-logout'),
  
  // 退出所有设备
  logoutAll: () => request.post('/manage/logout-all'),
  
  // 获取操作日志
  getLogs: (params) => request.get('/manage/logs', { params }),
  
  // 获取用户信息
  getProfile: () => request.get('/manage/profile')
}
