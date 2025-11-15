import request from '../utils/request'

// 登录
export const login = (data) => {
  return request.post('/auth/login', data)
}

// 注册
export const register = (data) => {
  return request.post('/auth/register', data)
}

// QQ登录
export const qqLogin = (data) => {
  return request.post('/auth/oauth/qq/login', data)
}

// 忘记密码
export const forgotPassword = (data) => {
  return request.post('/auth/forgotPassword', data)
}




