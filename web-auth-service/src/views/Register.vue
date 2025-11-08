<template>
  <div class="register-container">
    <div class="register-card">
      <h1 class="title">注册账号</h1>
      <p class="subtitle">{{ appName }}</p>

      <form @submit.prevent="handleRegister" class="register-form">
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
          <label>昵称</label>
          <input
            v-model="form.nickname"
            type="text"
            placeholder="请输入昵称"
            required
          />
        </div>

        <div class="form-group">
          <label>密码</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="8-20位字符"
            minlength="8"
            maxlength="20"
            required
          />
        </div>

        <div class="form-group">
          <label>确认密码</label>
          <input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
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

        <button type="submit" class="btn-register" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>

        <div class="links">
          <router-link to="/login">已有账号？立即登录</router-link>
        </div>
      </form>

      <!-- ✅ 预留固定高度空间，避免抖动 -->
      <div class="message-container">
        <p v-show="error" class="error-msg">{{ error }}</p>
        <p v-show="success" class="success-msg">{{ success }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

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
const error = ref('')
const success = ref('')
const appName = ref('')
const captchaImage = ref('')

const urlParams = new URLSearchParams(window.location.search)
const appId = urlParams.get('app_id') || 'blog'
const redirectUri = urlParams.get('redirect_uri') || 'http://localhost:3000/sso-callback'
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

async function handleRegister() {
  error.value = ''
  success.value = ''

  if (form.value.password !== form.value.confirmPassword) {
    error.value = '两次输入的密码不一致'
    return
  }

  loading.value = true

  try {
    const response = await axios.post('/api/auth/register', {
      email: form.value.email,
      nickname: form.value.nickname,
      password: form.value.password,
      captcha: form.value.captcha,
      captcha_id: form.value.captcha_id,
      app_id: appId
    })

    if (response.data.code === 0) {
      // ✅ 注册成功，跳转到登录页
      success.value = '注册成功！2秒后跳转到登录页面...'
      setTimeout(() => {
        router.push({
          path: '/login',
          query: {
            app_id: appId,
            redirect_uri: redirectUri,
            state: state
          }
        })
      }, 2000)
    } else {
      error.value = response.data.message
      refreshCaptcha()
    }
  } catch (err) {
    error.value = err.response?.data?.message || '注册失败，请重试'
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
}

.register-card {
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

.register-form {
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

.btn-register {
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

.btn-register:hover:not(:disabled) {
  transform: translateY(-2px);
}

.btn-register:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.links {
  text-align: center;
  font-size: 14px;
}

.links a {
  color: #667eea;
  text-decoration: none;
}

.links a:hover {
  text-decoration: underline;
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

.success-msg {
  padding: 12px;
  background: #efe;
  color: #3c3;
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

