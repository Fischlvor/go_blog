<template>
  <AuthLayout
    title="注册账号"
    :subtitle="`加入 ${appName}，开启您的旅程`"
    security-tip="注册即表示您同意我们的服务条款和隐私政策"
    :app-id="appId"
  >
    <template #icon>
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M16 21V19C16 17.9391 15.5786 16.9217 14.8284 16.1716C14.0783 15.4214 13.0609 15 12 15H5C3.93913 15 2.92172 15.4214 2.17157 16.1716C1.42143 16.9217 1 17.9391 1 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M8.5 11C10.7091 11 12.5 9.20914 12.5 7C12.5 4.79086 10.7091 3 8.5 3C6.29086 3 4.5 4.79086 4.5 7C4.5 9.20914 6.29086 11 8.5 11Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M23 21V19C22.9993 18.1137 22.7044 17.2528 22.1614 16.5523C21.6184 15.8519 20.8581 15.3516 20 15.13" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M16 3.13C16.8604 3.35031 17.623 3.85071 18.1676 4.55232C18.7122 5.25392 19.0078 6.11683 19.0078 7.005C19.0078 7.89318 18.7122 8.75608 18.1676 9.45769C17.623 10.1593 16.8604 10.6597 16 10.88" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </template>

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
          placeholder="请输入邮箱地址"
          required
          class="form-input"
          :disabled="registerSuccess"
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
          placeholder="请输入昵称"
          required
          class="form-input"
          :disabled="registerSuccess"
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
          :disabled="registerSuccess"
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
          :disabled="registerSuccess"
        />
        <p v-if="form.confirmPassword && form.password !== form.confirmPassword" class="error-text">
          两次密码不一致
        </p>
      </div>

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
            :disabled="registerSuccess"
          />
          <button
            type="button"
            class="btn-send-code"
            @click="openEmailCaptchaDialog"
            :disabled="sendingEmailCode || countdown > 0 || registerSuccess"
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

      <CaptchaDialog
        v-model="showEmailCaptchaDialog"
        :confirm-loading="sendingEmailCode"
        @confirm="onCaptchaConfirm"
      />

      <button
        type="submit"
        class="btn-register"
        :disabled="loading || registerSuccess || (form.password && form.password !== form.confirmPassword)"
      >
        <span v-if="!loading">注册</span>
        <span v-else class="loading-content">
          <span class="loading-dots"></span>
          注册中...
        </span>
      </button>

      <div class="links">
        <router-link
          :to="{ path: '/login', query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl } }"
          class="link-item"
        >
          已有账号？立即登录
        </router-link>
      </div>
    </form>
  </AuthLayout>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { sendEmailVerificationCode } from '@/api/base'
import { register } from '@/api/auth'
import AuthLayout from '@/components/AuthLayout.vue'
import CaptchaDialog from '@/components/CaptchaDialog.vue'

const router = useRouter()

const form = ref({
  email: '',
  nickname: '',
  password: '',
  confirmPassword: '',
  emailCode: ''
})

const loading = ref(false)
const sendingEmailCode = ref(false)
const showEmailCaptchaDialog = ref(false)
const countdown = ref(0)
const timer = ref(null)
const registerSuccess = ref(false)

const urlParams = new URLSearchParams(window.location.search)
const appId = urlParams.get('app_id') || 'blog'
const redirectUri = urlParams.get('redirect_uri') || 'http://localhost:3000/sso-callback'
const returnUrl = urlParams.get('return_url') || '/'

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
      scene: 'register',
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

async function handleRegister() {
  if (form.value.password !== form.value.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  loading.value = true
  registerSuccess.value = true
  try {
    const response = await register({
      email: form.value.email,
      nickname: form.value.nickname,
      password: form.value.password,
      verification_code: form.value.emailCode,
      app_id: appId
    })
    if (response.data.code === 0) {
      ElMessage.success('注册成功，3秒后跳转到登录页')
      setTimeout(() => {
        router.push({
          path: '/login',
          query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl }
        })
      }, 3000)
    } else {
      registerSuccess.value = false
      ElMessage.error(response.data.message || '注册失败')
    }
  } catch (err) {
    registerSuccess.value = false
    ElMessage.error(err.response?.data?.message || '注册失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>
