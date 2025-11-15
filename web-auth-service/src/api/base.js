import axios from 'axios'

// 获取验证码
export const getCaptcha = () => {
  return axios.get('/api/base/captcha')
}

// 获取QQ登录URL
export const getQQLoginURL = (appId, state) => {
  const params = { app_id: appId }
  if (state) {
    params.state = state
  }
  return axios.get('/api/base/qqLoginURL', { params })
}

// 发送邮箱验证码
export const sendEmailVerificationCode = (data) => {
  return axios.post('/api/auth/sendEmailVerificationCode', data)
}

