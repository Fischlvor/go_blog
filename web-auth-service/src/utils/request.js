import axios from 'axios'

// 创建 axios 实例
const service = axios.create({
  baseURL: import.meta.env.VITE_BASE_API || '/api',
  timeout: 10000
})

// 响应拦截器
service.interceptors.response.use(
  response => response,
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

export default service
