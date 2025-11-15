<template>
  <div class="dialog-overlay" @click.self="$emit('close')">
    <div class="dialog-content">
      <div class="dialog-header">
        <h2>忘记密码</h2>
        <button class="close-btn" @click="$emit('close')">×</button>
      </div>
      
      <form @submit.prevent="handleSubmit" class="forgot-password-form">
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
          <label>图片验证码</label>
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

        <div class="form-group">
          <label>邮箱验证码</label>
          <div class="verification-group">
            <input
              v-model="form.verification_code"
              placeholder="请输入邮箱验证码"
              maxlength="6"
              required
            />
            <button
              type="button"
              class="btn-send-code"
              @click="sendVerificationCode"
              :disabled="!canSendCode || sendingCode"
            >
              {{ sendingCode ? '发送中...' : '发送验证码' }}
            </button>
          </div>
        </div>

        <div class="form-group">
          <label>新密码</label>
          <input
            v-model="form.new_password"
            type="password"
            placeholder="请输入新密码（8-20位）"
            minlength="8"
            maxlength="20"
            required
          />
        </div>

        <div class="form-group">
          <label>确认密码</label>
          <input
            v-model="repeatPassword"
            type="password"
            placeholder="请再次输入新密码"
            required
          />
          <p v-if="form.new_password && form.new_password !== repeatPassword" class="error-text">
            两次密码不一致
          </p>
        </div>

        <div class="form-actions">
          <button type="button" class="btn-cancel" @click="$emit('close')">
            取消
          </button>
          <button type="submit" class="btn-submit" :disabled="loading || form.new_password !== repeatPassword">
            {{ loading ? '提交中...' : '确定' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getCaptcha, sendEmailVerificationCode } from '@/api/base'
import { forgotPassword } from '@/api/auth'

const props = defineProps({
  appId: {
    type: String,
    default: 'blog'
  },
  redirectUri: {
    type: String,
    default: ''
  },
  state: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close'])

const form = ref({
  email: '',
  captcha: '',
  captcha_id: '',
  verification_code: '',
  new_password: ''
})

const repeatPassword = ref('')
const loading = ref(false)
const sendingCode = ref(false)
const captchaImage = ref('')

const canSendCode = computed(() => {
  return form.value.email && form.value.captcha && form.value.captcha.length === 6
})

onMounted(() => {
  refreshCaptcha()
})

async function refreshCaptcha() {
  try {
    const response = await getCaptcha()
    if (response.data.code === 0) {
      form.value.captcha_id = response.data.data.captcha_id
      captchaImage.value = response.data.data.pic_path
      form.value.captcha = ''
    }
  } catch (err) {
    console.error('获取验证码失败:', err)
    ElMessage.error('获取验证码失败，请刷新页面')
  }
}

async function sendVerificationCode() {
  if (!canSendCode.value) {
    ElMessage.warning('请先输入邮箱和图片验证码')
    return
  }

  sendingCode.value = true

  try {
    const response = await sendEmailVerificationCode({
      email: form.value.email,
      captcha: form.value.captcha,
      captcha_id: form.value.captcha_id
    })

    if (response.data.code === 0) {
      ElMessage.success('验证码已发送，请查收邮箱')
    } else {
      ElMessage.error(response.data.message || '发送验证码失败')
      refreshCaptcha()
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '发送验证码失败')
    refreshCaptcha()
  } finally {
    sendingCode.value = false
  }
}

async function handleSubmit() {
  if (form.value.new_password !== repeatPassword.value) {
    ElMessage.error('两次密码不一致')
    return
  }

  if (form.value.new_password.length < 8 || form.value.new_password.length > 20) {
    ElMessage.error('密码长度应为8-20位')
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
      ElMessage.success('密码重置成功，请重新登录')
      setTimeout(() => {
        emit('close')
      }, 2000)
    } else {
      ElMessage.error(response.data.message || '密码重置失败')
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '密码重置失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog-content {
  background: white;
  border-radius: 16px;
  padding: 32px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.dialog-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 32px;
  color: #999;
  cursor: pointer;
  line-height: 1;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #333;
}

.forgot-password-form {
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

.verification-group {
  display: flex;
  gap: 12px;
}

.verification-group input {
  flex: 1;
}

.btn-send-code {
  padding: 12px 20px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  white-space: nowrap;
  transition: opacity 0.2s;
}

.btn-send-code:hover:not(:disabled) {
  opacity: 0.9;
}

.btn-send-code:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-text {
  color: #c33;
  font-size: 12px;
  margin: 4px 0 0 0;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;
}

.btn-cancel,
.btn-submit {
  flex: 1;
  padding: 14px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.btn-cancel {
  background: #f5f5f5;
  color: #666;
}

.btn-cancel:hover {
  background: #e8e8e8;
}

.btn-submit {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>

