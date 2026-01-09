<template>
  <AuthLayout
    title="找回密码"
    subtitle="输入您的邮箱，我们将发送验证码帮您重置密码"
    security-tip="您的密码将被安全加密存储"
    :app-id="appId"
  >
    <template #icon>
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M12 8V12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M12 16H12.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M3 12H21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" opacity="0.3"/>
        <path d="M12 3V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" opacity="0.3"/>
      </svg>
    </template>

    <form @submit.prevent="handleResetPassword" class="forgot-password-form">
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
            :disabled="sendingCode || countdown > 0"
          >
            <span v-if="countdown > 0">{{ countdown }}秒后重试</span>
            <span v-else-if="!sendingCode">发送验证码</span>
            <span v-else class="loading-content">
              <span class="loading-dots"></span>
              发送中...
            </span>
          </button>
        </div>
      </div>

      <div class="form-group">
        <label>
          <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M7 11V7C7 5.67392 7.52678 4.40215 8.46447 3.46447C9.40215 2.52678 10.6739 2 12 2C13.3261 2 14.5979 2.52678 15.5355 3.46447C16.4732 4.40215 17 5.67392 17 7V11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          新密码
        </label>
        <input
          v-model="form.new_password"
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
          确认新密码
        </label>
        <input
          v-model="repeatPassword"
          type="password"
          placeholder="请再次输入新密码"
          required
          class="form-input"
        />
        <p v-if="repeatPassword && form.new_password !== repeatPassword" class="error-text">
          两次密码不一致
        </p>
      </div>

      <CaptchaDialog
        v-model="showEmailCaptchaDialog"
        :confirm-loading="sendingCode"
        @confirm="onCaptchaConfirm"
      />

      <button
        type="submit"
        class="btn-submit"
        :disabled="loading || (repeatPassword && form.new_password !== repeatPassword)"
      >
        <span v-if="!loading">重置密码</span>
        <span v-else class="loading-content">
          <span class="loading-dots"></span>
          重置中...
        </span>
      </button>

      <div class="links">
        <router-link
          :to="{ path: '/login', query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl } }"
          class="link-item"
        >
          返回登录
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
import { forgotPassword } from '@/api/auth'
import AuthLayout from '@/components/AuthLayout.vue'
import CaptchaDialog from '@/components/CaptchaDialog.vue'

const router = useRouter()

const form = ref({
  email: '',
  verification_code: '',
  new_password: ''
})

const repeatPassword = ref('')
const loading = ref(false)
const sendingCode = ref(false)
const showEmailCaptchaDialog = ref(false)
const countdown = ref(0)
const timer = ref(null)

const urlParams = new URLSearchParams(window.location.search)
const appId = urlParams.get('app_id') || 'blog'
const redirectUri = urlParams.get('redirect_uri') || 'http://localhost:3000/sso-callback'
const returnUrl = urlParams.get('return_url') || '/'

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const canSendCode = computed(() => {
  return !!form.value.email && emailRegex.test(form.value.email)
})

async function openEmailCaptchaDialog() {
  if (!canSendCode.value) {
    ElMessage.warning('请先输入正确的邮箱地址')
    return
  }
  showEmailCaptchaDialog.value = true
}

async function onCaptchaConfirm(payload) {
  sendingCode.value = true
  showEmailCaptchaDialog.value = false
  try {
    const response = await sendEmailVerificationCode({
      email: form.value.email,
      scene: 'forgot_password',
      captcha: payload.captcha,
      captcha_id: payload.captcha_id
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
    sendingCode.value = false
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

async function handleResetPassword() {
  if (form.value.new_password !== repeatPassword.value) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  loading.value = true
  try {
    const response = await forgotPassword({
      email: form.value.email,
      verification_code: form.value.verification_code,
      new_password: form.value.new_password
    })
    if (response.data.code === 0) {
      ElMessage.success('密码重置成功！2秒后跳转到登录页面...')
      setTimeout(() => {
        router.push({
          path: '/login',
          query: { app_id: appId, redirect_uri: redirectUri, return_url: returnUrl }
        })
      }, 2000)
    } else {
      ElMessage.error(response.data.message || '密码重置失败')
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '密码重置失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>
