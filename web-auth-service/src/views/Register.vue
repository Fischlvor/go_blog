<template>
  <div class="register-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>

    <div class="register-card">
      <!-- Logo区域 -->
      <div class="logo-section">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M16 21V19C16 17.9391 15.5786 16.9217 14.8284 16.1716C14.0783 15.4214 13.0609 15 12 15H5C3.93913 15 2.92172 15.4214 2.17157 16.1716C1.42143 16.9217 1 17.9391 1 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M8.5 11C10.7091 11 12.5 9.20914 12.5 7C12.5 4.79086 10.7091 3 8.5 3C6.29086 3 4.5 4.79086 4.5 7C4.5 9.20914 6.29086 11 8.5 11Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M23 21V19C22.9993 18.1137 22.7044 17.2528 22.1614 16.5523C21.6184 15.8519 20.8581 15.3516 20 15.13" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M16 3.13C16.8604 3.35031 17.623 3.85071 18.1676 4.55232C18.7122 5.25392 19.0078 6.11683 19.0078 7.005C19.0078 7.89318 18.7122 8.75608 18.1676 9.45769C17.623 10.1593 16.8604 10.6597 16 10.88" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="logo-text">
          <h1 class="title">注册账号</h1>
          <p class="subtitle">加入 {{ appName }}，开启您的旅程</p>
        </div>
      </div>

      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 4H20C21.1 4 22 4.9 22 6V18C22 19.1 21.1 20 20 20H4C2.9 20 2 19.1 2 18V6C2 4.9 2.9 4 4 4Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M22 6L12 13L2 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            邮箱
          </label>
          <input
            v-model="form.email"
            type="email"
            placeholder="请输入您的邮箱地址"
            required
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M12 11C14.2091 11 16 9.20914 16 7C16 4.79086 14.2091 3 12 3C9.79086 3 8 4.79086 8 7C8 9.20914 9.79086 11 12 11Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            昵称
          </label>
          <input
            v-model="form.nickname"
            type="text"
            placeholder="请输入您的昵称"
            required
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M7 11V7C7 5.67392 7.52678 4.40215 8.46447 3.46447C9.40215 2.52678 10.6739 2 12 2C13.3261 2 14.5979 2.52678 15.5355 3.46447C16.4732 4.40215 17 5.67392 17 7V11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            密码
          </label>
          <input
            v-model="form.password"
            type="password"
            placeholder="8-20位字符"
            minlength="8"
            maxlength="20"
            required
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M7 11V7C7 5.67392 7.52678 4.40215 8.46447 3.46447C9.40215 2.52678 10.6739 2 12 2C13.3261 2 14.5979 2.52678 15.5355 3.46447C16.4732 4.40215 17 5.67392 17 7V11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            确认密码
          </label>
          <input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            required
            class="form-input"
          />
          <p v-if="form.confirmPassword && form.password !== form.confirmPassword" class="error-text">
            两次密码不一致
          </p>
        </div>

        <div class="form-group">
          <label>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M9 9H15M9 15H15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
            验证码
          </label>
          <div class="captcha-group">
            <input
              v-model="form.captcha"
              placeholder="请输入验证码"
              maxlength="6"
              required
              class="form-input"
              @focus="loadCaptchaIfNeeded"
            />
            <div class="captcha-wrapper" @click="handleCaptchaClick">
              <img
                v-if="captchaImage"
                :src="captchaImage"
                class="captcha-image"
                alt="验证码"
              />
              <div v-else class="captcha-placeholder">
                <span>点击加载</span>
              </div>
            </div>
          </div>
        </div>

        <button type="submit" class="btn-register" :disabled="loading || (form.password && form.password !== form.confirmPassword)">
          <span v-if="!loading">注册</span>
          <span v-else class="loading-content">
            <span class="loading-dots"></span>
            注册中...
          </span>
        </button>

        <div class="links">
          <router-link :to="{ path: '/login', query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl } }" class="link-item">
            已有账号？立即登录
          </router-link>
        </div>
      </form>

      <!-- 安全提示 -->
      <div class="security-tip">
        <svg class="tip-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 8V12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 16H12.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>注册即表示您同意我们的服务条款和隐私政策</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const router = useRouter()

const form = ref({
  email: '',
  nickname: '',
  password: '',
  confirmPassword: '',
  captcha: '',
  captcha_id: ''
})

const loading = ref(false)
const appName = ref('')
const captchaImage = ref('')

const urlParams = new URLSearchParams(window.location.search)
const appId = urlParams.get('app_id') || 'blog'
const redirectUri = urlParams.get('redirect_uri') || 'http://localhost:3000/sso-callback'
const returnUrl = urlParams.get('return_url') || '/'

onMounted(() => {
  appName.value = getAppName(appId)
  // 不再自动加载验证码，改为懒加载
})

function getAppName(appId) {
  const appNames = {
    'blog': 'Go博客系统'
  }
  return appNames[appId] || '应用'
}

// 懒加载验证码
const captchaLoaded = ref(false)

async function loadCaptchaIfNeeded() {
  if (!captchaLoaded.value) {
    await refreshCaptcha()
  }
}

async function handleCaptchaClick() {
  // 如果还没加载过，先加载；否则刷新
  await refreshCaptcha()
}

async function refreshCaptcha() {
  try {
    const response = await request.get('/base/captcha')
    if (response.data.code === 0) {
      form.value.captcha_id = response.data.data.captcha_id
      captchaImage.value = response.data.data.pic_path
      captchaLoaded.value = true
    }
  } catch (err) {
    console.error('获取验证码失败:', err)
    ElMessage.error('获取验证码失败，请重试')
  }
}

async function handleRegister() {
  if (form.value.password !== form.value.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  loading.value = true

  try {
    const response = await request.post('/auth/register', {
      email: form.value.email,
      nickname: form.value.nickname,
      password: form.value.password,
      captcha: form.value.captcha,
      captcha_id: form.value.captcha_id,
      app_id: appId
    })

    if (response.data.code === 0) {
      // ✅ 注册成功，跳转到登录页
      ElMessage.success('注册成功！2秒后跳转到登录页面...')
      setTimeout(() => {
        router.push({
          path: '/login',
          query: {
            app_id: appId,
            redirect_uri: redirectUri,
            return_url: returnUrl
          }
        })
      }, 2000)
    } else {
      ElMessage.error(response.data.message || '注册失败')
      refreshCaptcha()
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '注册失败，请重试')
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  overflow: hidden;
  box-sizing: border-box;
  width: 100%;
  max-width: 100vw;
}

/* 背景装饰 */
.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 0;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  animation: float 20s infinite ease-in-out;
}

.circle-1 {
  width: 300px;
  height: 300px;
  top: -100px;
  right: -100px;
  animation-delay: 0s;
}

.circle-2 {
  width: 200px;
  height: 200px;
  bottom: -50px;
  left: -50px;
  animation-delay: 5s;
}

.circle-3 {
  width: 150px;
  height: 150px;
  top: 50%;
  left: 10%;
  animation-delay: 10s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -30px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

.register-card {
  position: relative;
  z-index: 1;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(255, 255, 255, 0.5);
  padding: 32px 32px;
  width: 100%;
  max-width: min(480px, calc(100vw - 40px));
  min-width: 0;
  box-sizing: border-box;
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Logo区域 */
.logo-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.logo-icon {
  width: 48px;
  height: 48px;
  flex-shrink: 0;
  color: #667eea;
  animation: pulse 2s infinite;
}

.logo-text {
  flex: 1;
  text-align: left;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.title {
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 4px;
  letter-spacing: -0.5px;
  line-height: 1.2;
}

.subtitle {
  font-size: 13px;
  color: #6b7280;
  font-weight: 400;
  line-height: 1.3;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.input-icon {
  width: 18px;
  height: 18px;
  color: #9ca3af;
  flex-shrink: 0;
}

.form-input {
  padding: 12px 14px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: #fff;
  color: #111827;
  box-sizing: border-box;
  width: 100%;
  max-width: 100%;
}

.form-input::placeholder {
  color: #9ca3af;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.error-text {
  color: #ef4444;
  font-size: 12px;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.captcha-group {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  width: 100%;
  min-width: 0;
}

.captcha-group .form-input {
  flex: 1;
}

.captcha-wrapper {
  width: 130px;
  height: 48px;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid #e5e7eb;
  transition: all 0.3s;
  flex-shrink: 0;
}

.captcha-wrapper:hover {
  border-color: #667eea;
}

.captcha-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.captcha-placeholder {
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.captcha-placeholder span {
  color: #999;
  font-size: 12px;
  font-weight: 400;
}

.captcha-loading {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.btn-register {
  padding: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  position: relative;
  overflow: hidden;
}

.btn-register::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  transform: translate(-50%, -50%);
  transition: width 0.6s, height 0.6s;
}

.btn-register:hover:not(:disabled)::before {
  width: 300px;
  height: 300px;
}

.btn-register:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.5);
}

.btn-register:active:not(:disabled) {
  transform: translateY(0);
}

.btn-register:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.loading-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.loading-dots {
  display: inline-flex;
  gap: 4px;
}

.loading-dots::after {
  content: '...';
  animation: dots 1.5s steps(4, end) infinite;
}

@keyframes dots {
  0%, 20% {
    content: '.';
  }
  40% {
    content: '..';
  }
  60%, 100% {
    content: '...';
  }
}

.links {
  text-align: center;
  font-size: 14px;
}

.link-item {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
  position: relative;
}

.link-item::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: #667eea;
  transition: width 0.3s;
}

.link-item:hover::after {
  width: 100%;
}

.security-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 16px;
  padding: 10px;
  background: #f0f9ff;
  border-radius: 8px;
  font-size: 11px;
  color: #0369a1;
}

.tip-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  color: #0284c7;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .register-container {
    padding: 16px;
  }

  .register-card {
    padding: 32px 20px;
    border-radius: 20px;
    max-width: calc(100vw - 32px);
  }

  .title {
    font-size: 28px;
  }

  .logo-icon {
    width: 56px;
    height: 56px;
  }
}
</style>

