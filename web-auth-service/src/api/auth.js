import axios from 'axios'

// 登录
export const login = (data) => {
  return axios.post('/api/auth/login', data)
}

// 注册
export const register = (data) => {
  return axios.post('/api/auth/register', data)
}

// QQ登录
export const qqLogin = (data) => {
  return axios.post('/api/auth/oauth/qq/login', data)
}

// 忘记密码
export const forgotPassword = (data) => {
  return axios.post('/api/auth/forgotPassword', data)
}




