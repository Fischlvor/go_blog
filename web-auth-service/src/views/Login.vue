<template>
  <div class="login-container" :data-app-theme="appTheme">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>

    <div class="login-card">
      <!-- 错误提示 -->
      <div v-if="!hasValidAppId" class="error-section">
        <div class="error-icon">⚠️</div>
        <h2 class="error-title">访问错误</h2>
        <p class="error-message">请通过正确的链接访问登录页面</p>
        <div class="error-help">
          <p>如果您需要帮助，请联系系统管理员</p>
        </div>
      </div>

      <!-- 正常登录界面 -->
      <div v-else>
        <!-- Logo区域 -->
        <div class="logo-section">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <div class="logo-text">
            <h1 class="title">统一登录</h1>
            <p class="subtitle">欢迎使用 {{ appName }}</p>
          </div>
        </div>

      <!-- 登录方式切换 -->
      <div class="login-type-tabs">
        <div 
          class="tab-item" 
          :class="{ active: loginType === 'password' }"
          @click="loginType = 'password'"
        >
          密码登录
        </div>
        <div 
          class="tab-item" 
          :class="{ active: loginType === 'code' }"
          @click="loginType = 'code'"
        >
          验证码登录
        </div>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
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

        <!-- 密码登录方式 -->
        <template v-if="loginType === 'password'">
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
              placeholder="请输入您的密码"
              required
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label>
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M9 9H15M9 15H15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
              图片验证码
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
        </template>

        <!-- 验证码登录方式 -->
        <template v-else>
          <div class="form-group">
            <label>
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M21 10C21 17 12 23 12 23C12 23 3 17 3 10C3 7.61305 3.94821 5.32387 5.63604 3.63604C7.32387 1.94821 9.61305 1 12 1C14.3869 1 16.6761 1.94821 18.364 3.63604C20.0518 5.32387 21 7.61305 21 10Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M12 13C13.6569 13 15 11.6569 15 10C15 8.34315 13.6569 7 12 7C10.3431 7 9 8.34315 9 10C9 11.6569 10.3431 13 12 13Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              邮箱验证码
            </label>
            <div class="verification-group">
              <input
                v-model="form.emailCode"
                placeholder="请输入邮箱验证码"
                maxlength="6"
                required
                class="form-input"
              />
              <button
                type="button"
                class="btn-send-code"
                @click="openEmailCaptchaDialog"
                :disabled="sendingEmailCode"
              >
                <span v-if="!sendingEmailCode">发送验证码</span>
                <span v-else class="loading-content">
                  <span class="loading-dots"></span>
                  发送中...
                </span>
              </button>
            </div>
          </div>
        </template>

        <button type="submit" class="btn-login" :disabled="loading">
          <span v-if="!loading">登录</span>
          <span v-else class="loading-content">
            <span class="loading-dots"></span>
            登录中...
          </span>
        </button>

        <div class="links">
          <router-link :to="{ path: '/register', query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl, state: state } }" class="link-item">
            注册账号
          </router-link>
          <span class="divider-dot">·</span>
          <router-link :to="{ path: '/forgot-password', query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl, state: state } }" class="link-item">
            忘记密码?
          </router-link>
        </div>

        <div class="divider">
          <span class="divider-text">或使用第三方登录</span>
        </div>

        <div class="oauth-login">
          <div class="oauth-icon-wrapper" @click="qqLogin" title="QQ登录">
            <img 
              src="https://image.hsk423.cn/blog/8d17d9baa96c667b6dfbf7f2dfff42a2-20251111123335.jpg" 
              alt="QQ" 
              class="oauth-icon-img"
            />
          </div>
        </div>
      </form>

      <CaptchaDialog
        v-model="showEmailCaptchaDialog"
        :confirm-loading="sendingEmailCode"
        @confirm="onCaptchaConfirm"
      />

        <!-- 安全提示 -->
        <div class="security-tip">
          <svg class="tip-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M12 8V12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M12 16H12.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span>您的登录信息已加密传输，请放心使用</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCaptcha, getQQLoginURL, sendEmailVerificationCode } from '@/api/base'
import { login } from '@/api/auth'
import { getDeviceId, getBrowserName } from '@/utils/device'
import { encodeState, generateNonce } from '@/utils/state'
import { computed } from 'vue'
import CaptchaDialog from '@/components/CaptchaDialog.vue'

const loginType = ref('password') // 'password' 或 'code'

const form = ref({
  email: '',
  password: '',
  captcha: '',
  captcha_id: '',
  emailCode: ''
})

const loading = ref(false)
const sendingEmailCode = ref(false)
const appName = ref('')
const captchaImage = ref('')
const showEmailCaptchaDialog = ref(false)
// 使用通用验证码弹窗组件，不再在此维护图片验证码状态

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const canSendEmailCode = computed(() => {
  return !!form.value.email && emailRegex.test(form.value.email)
})

// 从URL获取参数
const urlParams = new URLSearchParams(window.location.search)
const appId = urlParams.get('app_id')

// 验证app_id是否存在
const hasValidAppId = ref(!!appId)
if (!appId) {
  ElMessage.error('访问链接无效，请联系管理员')
  hasValidAppId.value = false
}

// 根据app_id设置默认的redirect_uri
let defaultRedirectUri = 'http://localhost:3000/sso-callback'
if (appId === 'manage') {
  defaultRedirectUri = 'http://localhost:3001/manage'
}

const redirectUri = urlParams.get('redirect_uri') || defaultRedirectUri
const returnUrl = urlParams.get('return_url') || '/'
const state = urlParams.get('state') || ''

// 根据app_id设置主题
const appTheme = ref('default')

onMounted(() => {
  appName.value = getAppName(appId)
  appTheme.value = getAppTheme(appId)
  // 不再自动加载验证码，改为懒加载
  
  // 调试信息
  console.log('登录页面参数:', {
    appId,
    redirectUri,
    returnUrl,
    theme: appTheme.value
  })
})

function getAppName(appId) {
  const appNames = {
    'blog': 'Go博客系统',
    'manage': 'SSO管理中心',
    'mcp': 'MCP文档系统'
  }
  return appNames[appId] || '应用'
}

function getAppTheme(appId) {
  const themes = {
    'blog': 'blue',
    'mcp': 'green',
    'manage': 'purple'
  }
  return themes[appId] || 'purple'
}

const router = useRouter()
const route = useRoute()

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
    const response = await getCaptcha()
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

async function openEmailCaptchaDialog() {
  if (!canSendEmailCode.value) {
    ElMessage.warning('请先输入正确的邮箱地址')
    return
  }
  showEmailCaptchaDialog.value = true
}

async function onCaptchaConfirm(payload) {
  sendingEmailCode.value = true

  try {
    const response = await sendEmailVerificationCode({
      email: form.value.email,
      captcha_id: payload.captcha_id,
      captcha: payload.captcha
    })

    if (response.data.code === 0) {
      ElMessage.success('验证码已发送，请查收邮箱')
      showEmailCaptchaDialog.value = false
    } else {
      ElMessage.error(response.data.message || '发送验证码失败')
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '发送验证码失败')
  } finally {
    sendingEmailCode.value = false
  }
}

async function handleLogin() {
  // 验证app_id是否有效
  if (!appId) {
    ElMessage.error('访问链接无效，无法登录')
    return
  }

  loading.value = true

  try {
    // 生成state参数
    const stateData = {
      nonce: generateNonce(),
      app_id: appId,
      device_id: getDeviceId(),
      redirect_uri: redirectUri,
      return_url: returnUrl
    }
    const encodedState = encodeState(stateData)

    let response
    if (loginType.value === 'password') {
      // 密码登录
      response = await login({
        email: form.value.email,
        password: form.value.password,
        captcha: form.value.captcha,
        captcha_id: form.value.captcha_id,
        state: encodedState,
        app_id: appId,  // 直接传递app_id参数
        device_name: getBrowserName(),
        device_type: 'web'
      })
    } else {
      // 验证码登录
      response = await login({
        email: form.value.email,
        verification_code: form.value.emailCode,
        state: encodedState,
        app_id: appId,  // 直接传递app_id参数
        device_name: getBrowserName(),
        device_type: 'web'
      })
    }

    if (response.data.code === 0) {
      const data = response.data.data
      if (appId === 'manage') {
        // 管理后台登录：存储Token并跳转
        if (data.access_token) {
          localStorage.setItem('access_token', data.access_token)
          if (data.refresh_token) {
            localStorage.setItem('refresh_token', data.refresh_token)
          }
          // 使用路由跳转而不是window.location.href
          router.push('/manage')
        } else {
          ElMessage.error('登录失败：未获取到访问令牌')
        }
      } else {
        // ✅ OAuth 2.0: 回调时携带code和return_url
        let callbackUrl = `${data.redirect_uri}?code=${data.code}`
        if (data.return_url) {
          callbackUrl += `&return_url=${encodeURIComponent(data.return_url)}`
        }
        window.location.href = callbackUrl
      }
    } else {
      ElMessage.error(response.data.message || '登录失败')
      if (loginType.value === 'password') {
        refreshCaptcha()
      }
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '登录失败，请重试')
    if (loginType.value === 'password') {
      refreshCaptcha()
    }
  } finally {
    loading.value = false
  }
}

// getDeviceId 和 getBrowserName 已从 @/utils/device 导入

async function qqLogin() {
  // 验证app_id是否有效
  if (!appId) {
    ElMessage.error('访问链接无效，无法登录')
    return
  }

  try {
    // 生成state参数
    const stateData = {
      nonce: generateNonce(),
      app_id: appId,
      device_id: getDeviceId(),
      redirect_uri: redirectUri,
      return_url: returnUrl
    }
    const stateParam = encodeState(stateData)
    
    const response = await getQQLoginURL(appId, stateParam)
    if (response.data.code === 0) {
      // 跳转到QQ登录页面
      window.location.href = response.data.data.url
    } else {
      ElMessage.error(response.data.message || '获取QQ登录链接失败')
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '获取QQ登录链接失败')
  }
}
</script>

<style scoped>
.login-container {
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  background-size: 200% 200%;
  animation: gradientShift 15s ease infinite;
}

/* Blog 主题背景 - 蓝色主导+青色点缀 */
.login-container[data-app-theme="blue"] {
  background: linear-gradient(135deg, #bfdbfe 0%, #a5f3fc 50%, #bfdbfe 100%);
  background-size: 200% 200%;
  animation: gradientShift 15s ease infinite;
}

/* MCP 主题背景 - 绿色主导+黄绿点缀 */
.login-container[data-app-theme="green"] {
  background: linear-gradient(135deg, #a7f3d0 0%, #fde68a 50%, #a7f3d0 100%);
  background-size: 200% 200%;
  animation: gradientShift 15s ease infinite;
}

/* 默认/Manage 主题背景 - 紫色主导+粉色点缀 */
.login-container[data-app-theme="purple"] {
  background: linear-gradient(135deg, #ddd6fe 0%, #fbcfe8 50%, #ddd6fe 100%);
  background-size: 200% 200%;
  animation: gradientShift 15s ease infinite;
}

@keyframes gradientShift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
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
  animation: float 20s infinite ease-in-out;
}

/* Blog 主题 - 浅蓝色 */
.login-container[data-app-theme="blue"] .circle {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.15) 0%, rgba(37, 99, 235, 0.15) 100%);
}

/* MCP 主题 - 绿色 */
.login-container[data-app-theme="green"] .circle {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.15) 0%, rgba(5, 150, 105, 0.15) 100%);
}

/* 默认/Manage 主题 - 紫色 */
.login-container[data-app-theme="purple"] .circle {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
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

.login-card {
  position: relative;
  z-index: 1;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(255, 255, 255, 0.5);
  padding: 32px 32px;
  width: 100%;
  max-width: 440px;
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

/* Blog 主题 Logo */
.login-container[data-app-theme="blue"] .logo-icon {
  color: #3b82f6;
}

/* MCP 主题 Logo */
.login-container[data-app-theme="green"] .logo-icon {
  color: #10b981;
}

/* 默认/Manage 主题 Logo */
.login-container[data-app-theme="purple"] .logo-icon {
  color: #667eea;
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

/* Blog 主题标题 */
.login-container[data-app-theme="blue"] .title {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* MCP 主题标题 */
.login-container[data-app-theme="green"] .title {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 默认/Manage 主题标题 */
.login-container[data-app-theme="purple"] .title {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 13px;
  color: #6b7280;
  font-weight: 400;
  line-height: 1.3;
}

.login-type-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  background: #f3f4f6;
  border-radius: 10px;
  padding: 4px;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 10px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.3s;
}

.tab-item:hover {
  color: #374151;
}

.tab-item.active {
  background: #fff;
  color: #667eea;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* Blog 主题激活标签 */
.login-container[data-app-theme="blue"] .tab-item.active {
  color: #3b82f6;
}

/* MCP 主题激活标签 */
.login-container[data-app-theme="green"] .tab-item.active {
  color: #10b981;
}

/* 默认/Manage 主题激活标签 */
.login-container[data-app-theme="purple"] .tab-item.active {
  color: #667eea;
}

.login-form {
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

.captcha-group {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.verification-group {
  display: flex;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.verification-group .form-input {
  flex: 1;
}

.btn-send-code {
  padding: 12px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  flex-shrink: 0;
}

/* Blog 主题发送按钮 */
.login-container[data-app-theme="blue"] .btn-send-code {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
}

/* MCP 主题发送按钮 */
.login-container[data-app-theme="green"] .btn-send-code {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

/* 默认/Manage 主题发送按钮 */
.login-container[data-app-theme="purple"] .btn-send-code {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.btn-send-code:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

/* Blog 主题发送按钮悬停 */
.login-container[data-app-theme="blue"] .btn-send-code:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

/* MCP 主题发送按钮悬停 */
.login-container[data-app-theme="green"] .btn-send-code:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

/* 默认/Manage 主题发送按钮悬停 */
.login-container[data-app-theme="purple"] .btn-send-code:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-send-code:active:not(:disabled) {
  transform: translateY(0);
}

.btn-send-code:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.email-captcha-dialog .email-captcha-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 8px 0;
}

.email-captcha-image {
  width: 200px;
  height: 64px;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid #e5e7eb;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
}

.email-captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.email-captcha-tip {
  font-size: 12px;
  color: #6b7280;
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

.btn-login {
  padding: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

/* Blog 主题登录按钮 */
.login-container[data-app-theme="blue"] .btn-login {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
}

/* MCP 主题登录按钮 */
.login-container[data-app-theme="green"] .btn-login {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

/* 默认/Manage 主题登录按钮 */
.login-container[data-app-theme="purple"] .btn-login {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.btn-login::before {
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

.btn-login:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.5);
}

/* Blog 主题登录按钮悬停 */
.login-container[data-app-theme="blue"] .btn-login:hover:not(:disabled) {
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.5);
}

/* MCP 主题登录按钮悬停 */
.login-container[data-app-theme="green"] .btn-login:hover:not(:disabled) {
  box-shadow: 0 8px 20px rgba(16, 185, 129, 0.5);
}

/* 默认/Manage 主题登录按钮悬停 */
.login-container[data-app-theme="purple"] .btn-login:hover:not(:disabled) {
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.5);
}

.btn-login:active:not(:disabled) {
  transform: translateY(0);
}

.btn-login:disabled {
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
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 12px;
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

.divider-dot {
  color: #d1d5db;
  font-weight: 300;
}

.divider {
  position: relative;
  text-align: center;
  margin: 10px 0;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(to right, transparent, #e5e7eb, transparent);
}

.divider-text {
  position: relative;
  padding: 0 20px;
  background: rgba(255, 255, 255, 0.95);
  color: #9ca3af;
  font-size: 13px;
}

.oauth-login {
  display: flex;
  justify-content: center;
  align-items: center;
}

.oauth-icon-wrapper {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 50%;
  padding: 8px;
}

.oauth-icon-wrapper:hover {
  background: #f3f4f6;
  transform: scale(1.1);
}

.oauth-icon-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: 50%;
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

/* 错误提示样式 */
.error-section {
  text-align: center;
  padding: 40px 20px;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-title {
  font-size: 24px;
  font-weight: 600;
  color: #dc2626;
  margin-bottom: 12px;
}

.error-message {
  font-size: 16px;
  color: #6b7280;
  margin-bottom: 24px;
  line-height: 1.5;
}

.error-help {
  background: #fef3c7;
  border: 1px solid #f59e0b;
  border-radius: 8px;
  padding: 16px;
  text-align: left;
}

.error-help p {
  font-weight: 600;
  color: #92400e;
  margin-bottom: 8px;
}

.error-help ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.error-help li {
  color: #78350f;
  margin-bottom: 4px;
}

.error-help code {
  background: #fbbf24;
  color: #78350f;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .login-container {
    padding: 16px;
  }

  .login-card {
    padding: 32px 20px;
    border-radius: 20px;
    max-width: 100%;
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

