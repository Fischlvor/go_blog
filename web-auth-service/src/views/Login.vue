<template>
  <div class="login-container">
    <div class="login-card">
      <h1 class="title">统一登录</h1>
      <p class="subtitle">{{ appName }}</p>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>邮箱</label>
          <input
            v-model="form.email"
            type="email"
            placeholder="请输入邮箱"
            required
          />
        </div>

        <div class="form-group">
          <label>密码</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            required
          />
        </div>

        <div class="form-group">
          <label>验证码</label>
          <div class="captcha-group">
            <input
              v-model="form.captcha"
              placeholder="验证码"
              maxlength="6"
              required
            />
            <img
              v-if="captchaImage"
              :src="captchaImage"
              class="captcha-image"
              @click="refreshCaptcha"
              alt="验证码"
            />
            <div v-else class="captcha-loading" @click="refreshCaptcha">
              加载中...
            </div>
          </div>
        </div>

        <button type="submit" class="btn-login" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>

        <div class="links">
          <router-link :to="{ path: '/register', query: { app_id: appId, redirect_uri: redirectUri, state: state } }">注册账号</router-link>
          <a href="#" @click.prevent="forgotPassword">忘记密码?</a>
        </div>

        <div class="divider">
          <span>或使用第三方登录</span>
        </div>

        <div class="oauth-buttons">
          <button type="button" class="btn-qq" @click="qqLogin">
            <span>QQ登录</span>
          </button>
        </div>
      </form>

      <!-- ✅ 预留固定高度空间，避免抖动 -->
      <div class="message-container">
        <p v-show="error" class="error-msg">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const form = ref({
  email: '',
  password: '',
  captcha: '',
  captcha_id: ''
})

const loading = ref(false)
const error = ref('')
const appName = ref('')
const captchaImage = ref('')

// 从URL获取参数
const urlParams = new URLSearchParams(window.location.search)
const appId = urlParams.get('app_id') || 'blog'
const redirectUri = urlParams.get('redirect_uri') || 'http://localhost:3000/callback'
const state = urlParams.get('state') || ''

onMounted(() => {
  appName.value = getAppName(appId)
  refreshCaptcha()
})

function getAppName(appId) {
  const appNames = {
    'blog': 'Go博客系统'
  }
  return appNames[appId] || '应用'
}

async function refreshCaptcha() {
  try {
    const response = await axios.get('/api/base/captcha')
    if (response.data.code === 0) {
      form.value.captcha_id = response.data.data.captcha_id
      captchaImage.value = response.data.data.pic_path
    }
  } catch (err) {
    console.error('获取验证码失败:', err)
    error.value = '获取验证码失败，请刷新页面'
  }
}

async function handleLogin() {
  error.value = ''
  loading.value = true

  try {
    const response = await axios.post('/api/auth/login', {
      email: form.value.email,
      password: form.value.password,
      captcha: form.value.captcha,
      captcha_id: form.value.captcha_id,
      app_id: appId,
      redirect_uri: redirectUri,
      device_id: getDeviceId(),
      device_name: getBrowserName(),
      device_type: 'web'
    })

    if (response.data.code === 0) {
      const data = response.data.data
      // ✅ 标准OAuth 2.0: 返回授权码
      window.location.href = `${data.redirect_uri}?` +
        `code=${data.code}&` +
        `state=${state}`
    } else {
      error.value = response.data.message
      refreshCaptcha()
    }
  } catch (err) {
    error.value = err.response?.data?.message || '登录失败，请重试'
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

function getDeviceId() {
  let deviceId = localStorage.getItem('device_id')
  if (!deviceId) {
    deviceId = 'web_' + Math.random().toString(36).substring(2, 15)
    localStorage.setItem('device_id', deviceId)
  }
  return deviceId
}

function getBrowserName() {
  const ua = navigator.userAgent
  if (ua.indexOf('Chrome') > -1) return 'Chrome'
  if (ua.indexOf('Firefox') > -1) return 'Firefox'
  if (ua.indexOf('Safari') > -1) return 'Safari'
  return 'Browser'
}

function qqLogin() {
  alert('QQ登录功能开发中...')
}

function forgotPassword() {
  alert('忘记密码功能开发中...')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 420px;
}

.title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  text-align: center;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 14px;
  color: #999;
  text-align: center;
  margin-bottom: 32px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #555;
}

.form-group input {
  padding: 12px 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.captcha-group {
  display: flex;
  gap: 12px;
}

.captcha-group input {
  flex: 1;
}

.captcha-image {
  width: 120px;
  height: 44px;
  border-radius: 8px;
  cursor: pointer;
  user-select: none;
  border: 1px solid #ddd;
  object-fit: contain;
}

.captcha-loading {
  width: 120px;
  height: 44px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  cursor: pointer;
  user-select: none;
}

.btn-login {
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.btn-login:hover:not(:disabled) {
  transform: translateY(-2px);
}

.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.links {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.links a {
  color: #667eea;
  text-decoration: none;
}

.links a:hover {
  text-decoration: underline;
}

.divider {
  position: relative;
  text-align: center;
  margin: 20px 0;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: #ddd;
}

.divider span {
  position: relative;
  padding: 0 16px;
  background: white;
  color: #999;
  font-size: 12px;
}

.oauth-buttons {
  display: flex;
  gap: 12px;
}

.btn-qq {
  flex: 1;
  padding: 12px;
  background: #12b7f5;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-qq:hover {
  opacity: 0.9;
}

.message-container {
  min-height: 60px;
  margin-top: 16px;
}

.error-msg {
  padding: 12px;
  background: #fee;
  color: #c33;
  border-radius: 8px;
  font-size: 14px;
  text-align: center;
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

