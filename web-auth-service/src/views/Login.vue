<template>
  <AuthLayout
    title="统一登录"
    :subtitle="`欢迎使用 ${appName}`"
    security-tip="您的信息已加密传输，请放心使用"
    :app-id="appId"
  >
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

      <!-- 密码登录 -->
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
            <div class="captcha-wrapper" @click="refreshCaptcha">
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

      <!-- 验证码登录 -->
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
              v-model="form.verification_code"
              placeholder="请输入邮箱验证码"
              maxlength="6"
              required
              class="form-input"
            />
            <button
              type="button"
              class="btn-send-code"
              @click="openEmailCaptchaDialog"
              :disabled="sendingEmailCode || countdown > 0"
            >
              <span v-if="countdown > 0">{{ countdown }}秒后重试</span>
              <span v-else-if="!sendingEmailCode">发送验证码</span>
              <span v-else class="loading-content">
                <span class="loading-dots"></span>
                发送中...
              </span>
            </button>
          </div>
        </div>
      </template>

      <CaptchaDialog
        v-model="showEmailCaptchaDialog"
        :confirm-loading="sendingEmailCode"
        @confirm="onCaptchaConfirm"
      />

      <button type="submit" class="btn-login" :disabled="loading || !canLogin">
        <span v-if="!loading">登录</span>
        <span v-else class="loading-content">
          <span class="loading-dots"></span>
          登录中...
        </span>
      </button>

      <div class="links">
        <router-link
          :to="{ path: '/register', query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl } }"
          class="link-item"
        >
          注册账号
        </router-link>
        <span class="divider-dot">·</span>
        <router-link
          :to="{ path: '/forgot-password', query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl } }"
          class="link-item"
        >
          忘记密码
        </router-link>
      </div>

      <!-- 第三方登录 -->
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
  </AuthLayout>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCaptcha, sendEmailVerificationCode, getQQLoginURL } from '@/api/base'
import { login } from '@/api/auth'
import { encodeState, generateNonce } from '@/utils/state'
import { getDeviceId, getBrowserName } from '@/utils/device'
import AuthLayout from '@/components/AuthLayout.vue'
import CaptchaDialog from '@/components/CaptchaDialog.vue'

const router = useRouter()

const form = ref({
  email: '',
  password: '',
  verification_code: '',
  captcha: '',
  captcha_id: ''
})

const loginType = ref('password')
const loading = ref(false)
const sendingEmailCode = ref(false)
const showEmailCaptchaDialog = ref(false)
const captchaImage = ref('')
const captchaLoaded = ref(false)
const countdown = ref(0)
const timer = ref(null)

const urlParams = new URLSearchParams(window.location.search)
const appId = urlParams.get('app_id') || 'blog'
const redirectUri = urlParams.get('redirect_uri') || 'http://localhost:3000/sso-callback'
const returnUrl = urlParams.get('return_url') || '/'
const state = urlParams.get('state') || ''

const appName = computed(() => {
  const appNames = {
    'blog': 'Go博客系统',
    'manage': 'SSO管理中心',
    'mcp': 'MCP文档系统'
  }
  return appNames[appId] || '应用'
})

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const canSendEmailCode = computed(() => {
  return !!form.value.email && emailRegex.test(form.value.email)
})

const canLogin = computed(() => {
  if (!form.value.email || !emailRegex.test(form.value.email)) {
    return false
  }
  if (loginType.value === 'password') {
    return !!form.value.password && !!form.value.captcha
  } else {
    return !!form.value.verification_code
  }
})

onMounted(() => {
  // 懒加载验证码
})

async function loadCaptchaIfNeeded() {
  if (!captchaLoaded.value) {
    await refreshCaptcha()
  }
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
  showEmailCaptchaDialog.value = false
  try {
    const response = await sendEmailVerificationCode({
      email: form.value.email,
      scene: 'login',
      captcha_id: payload.captcha_id,
      captcha: payload.captcha
    })
    if (response.data.code === 0) {
      ElMessage.success('验证码已发送，请查收邮箱')
      startCountdown(60)
    } else if (response.data.code === 1013) {
      const seconds = response.data.remaining_seconds || 60
      startCountdown(seconds)
      ElMessage.warning(response.data.message || `请等待 ${seconds} 秒后再试`)
    } else {
      ElMessage.error(response.data.message || '发送验证码失败')
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '发送验证码失败')
  } finally {
    sendingEmailCode.value = false
  }
}

function startCountdown(seconds = 60) {
  countdown.value = seconds
  if (timer.value) clearInterval(timer.value)
  timer.value = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer.value)
      timer.value = null
    }
  }, 1000)
}

onUnmounted(() => {
  if (timer.value) {
    clearInterval(timer.value)
  }
})

async function handleLogin() {
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

    const loginData = {
      email: form.value.email,
      app_id: appId,
      state: encodedState,
      device_name: getBrowserName(),
      device_type: 'web'
    }

    if (loginType.value === 'password') {
      loginData.password = form.value.password
      loginData.captcha = form.value.captcha
      loginData.captcha_id = form.value.captcha_id
    } else {
      loginData.verification_code = form.value.verification_code
    }

    const response = await login(loginData)

    if (response.data.code === 0) {
      const data = response.data.data
      // OAuth 2.0 流程：构造回调URL
      let callbackUrl = `${data.redirect_uri}?code=${data.code}`
      if (data.return_url) {
        callbackUrl += `&return_url=${encodeURIComponent(data.return_url)}`
      }
      window.location.href = callbackUrl
    } else {
      ElMessage.error(response.data.message || '登录失败')
      if (loginType.value === 'password') {
        await refreshCaptcha()
        form.value.captcha = ''
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

async function qqLogin() {
  try {
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
.login-type-tabs {
  display: flex;
  background: #f3f4f6;
  border-radius: 10px;
  padding: 4px;
  margin-bottom: 20px;
}

.tab-item {
  flex: 1;
  padding: 10px;
  text-align: center;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.3s;
}

.tab-item.active {
  background: white;
  color: #667eea;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
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
  font-size: 12px;
}

.oauth-login {
  display: flex;
  justify-content: center;
  gap: 20px;
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
</style>
