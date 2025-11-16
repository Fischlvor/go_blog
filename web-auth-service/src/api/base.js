import request from '../utils/request'

// 获取验证码
export const getCaptcha = () => {
  return request.get('/base/captcha')
}

// 获取QQ登录URL
export const getQQLoginURL = (appId, state) => {
  const params = { app_id: appId }
  if (state) {
    params.state = state
  }
  return request.get('/base/qqLoginURL', { params })
}

// 发送邮箱验证码
export const sendEmailVerificationCode = (data) => {
  return request.post('/auth/sendEmailVerificationCode', data)
}
