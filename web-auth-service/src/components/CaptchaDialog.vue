<template>
  <el-dialog
    v-model="innerVisible"
    width="320px"
    :close-on-click-modal="false"
    :lock-scroll="true"
    append-to-body
    custom-class="captcha-dialog"
  >
    <template #header>
      <header class="dialog-header">
        <h3 class="dialog-title">{{ title }}</h3>
      </header>
    </template>

    <section class="captcha-content">
      <div class="captcha-preview" @click="loadCaptcha">
        <transition name="fade" mode="out-in">
          <img
            v-if="captchaImage"
            :key="captchaImage"
            :src="captchaImage"
            alt="图片验证码"
          />
          <div v-else key="loading" class="captcha-loading">
            <el-icon class="loading-icon" :size="20">
              <Loading />
            </el-icon>
            <span>验证码加载中…</span>
          </div>
        </transition>
        <el-button
          link
          class="refresh-btn"
          :loading="loading"
          @click.stop="loadCaptcha"
        >
          <el-icon><RefreshRight /></el-icon>
        </el-button>
      </div>

      <div class="code-input" @paste.prevent="onPaste">
        <input
          v-for="(ch, idx) in codeLength"
          :key="idx"
          :ref="(el) => setInputRef(el, idx)"
          class="code-box"
          type="text"
          inputmode="numeric"
          autocomplete="one-time-code"
          maxlength="1"
          :value="codeChars[idx]"
          @input="onInput(idx, $event)"
          @keydown="onKeydown(idx, $event)"
        />
      </div>
    </section>

  </el-dialog>
</template>

<script setup>
import { ref, watch, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { RefreshRight, Loading } from '@element-plus/icons-vue'
import { getCaptcha } from '@/api/base'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, default: '请输入图片验证码' },
  confirmLoading: { type: Boolean, default: false }
})

const emits = defineEmits(['update:modelValue', 'confirm'])

const innerVisible = ref(props.modelValue)
const loading = computed(() => props.confirmLoading)
watch(() => props.modelValue, (val) => {
  innerVisible.value = val
  if (val) loadCaptcha()
})

const captchaId = ref('')
const captchaImage = ref('')
const captchaValue = ref('')
const codeLength = 6
const codeChars = ref(Array.from({ length: codeLength }, () => ''))
const inputRefs = ref([])

async function loadCaptcha() {
  if (loading.value) return
  try {
    const res = await getCaptcha()
    if (res.data.code === 0) {
      captchaId.value = res.data.data.captcha_id
      captchaImage.value = res.data.data.pic_path
      resetCode()
      await nextTick(() => {
        focusCell(0)
      })
    }
  } catch (e) {
    ElMessage.error('获取验证码失败，请稍后重试')
  }
}

function confirm() {
  if (loading.value) return
  const value = captchaValue.value.trim()
  if (!value || value.length !== 6) {
    ElMessage.warning('请填写图片验证码')
    return
  }
  emits('confirm', { captcha: value, captcha_id: captchaId.value })
}

function close() {
  innerVisible.value = false
}

function resetCode() {
  captchaValue.value = ''
  codeChars.value = Array.from({ length: codeLength }, () => '')
  inputRefs.value = []
}

function focusCell(index) {
  const target = inputRefs.value[index]
  if (target) {
    target.focus()
    target.select && target.select()
  }
}

function setInputRef(el, index) {
  if (el) {
    inputRefs.value[index] = el
  } else {
    inputRefs.value[index] = undefined
  }
}

function normalizeChar(char) {
  const m = (char || '').match(/[0-9a-zA-Z]/)
  return m ? m[0] : ''
}

function updateCaptchaValueAndMaybeSubmit() {
  captchaValue.value = codeChars.value.join('')
  if (captchaValue.value.length === codeLength) {
    confirm()
  }
}

function onInput(index, evt) {
  const target = evt.target
  const value = normalizeChar(target.value)
  codeChars.value[index] = value
  target.value = value
  if (value && index < codeLength - 1) {
    focusCell(index + 1)
  }
  updateCaptchaValueAndMaybeSubmit()
}

function onKeydown(index, evt) {
  const key = evt.key
  if (key === 'Backspace') {
    if (codeChars.value[index]) {
      codeChars.value[index] = ''
      updateCaptchaValueAndMaybeSubmit()
    } else if (index > 0) {
      focusCell(index - 1)
      if (inputRefs.value[index - 1]) {
        inputRefs.value[index - 1].value = ''
      }
      codeChars.value[index - 1] = ''
      updateCaptchaValueAndMaybeSubmit()
    }
    evt.preventDefault()
  } else if (key === 'ArrowLeft' && index > 0) {
    focusCell(index - 1)
    evt.preventDefault()
  } else if (key === 'ArrowRight' && index < codeLength - 1) {
    focusCell(index + 1)
    evt.preventDefault()
  }
}

function onPaste(evt) {
  const text = (evt.clipboardData && evt.clipboardData.getData('text')) || ''
  if (!text) return
  const filtered = text.replace(/[^0-9a-zA-Z]/g, '').slice(0, codeLength)
  if (!filtered) return
  for (let i = 0; i < codeLength; i++) {
    codeChars.value[i] = filtered[i] || ''
  }
  nextTick(() => {
    const focusIndex = Math.min(filtered.length, codeLength - 1)
    focusCell(focusIndex)
    updateCaptchaValueAndMaybeSubmit()
  })
}

watch(innerVisible, (val) => {
  emits('update:modelValue', val)
  if (!val) {
    resetCode()
  }
})

watch(
  () => props.confirmLoading,
  (val) => {
    if (!val && innerVisible.value && !captchaValue.value) {
      nextTick(() => focusCell(0))
    }
  }
)
</script>

<style scoped>
.captcha-dialog .el-dialog__header {
  margin: 0;
  padding: 16px 20px 0;
}
.captcha-dialog .el-dialog__body {
  padding: 12px 20px 16px;
}
.dialog-header {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.dialog-title {
  font-size: 18px;
  font-weight: 500;
  margin: 0;
  color: #1e293b;
}
.dialog-subtitle {
  font-size: 12px;
  color: #94a3b8;
  margin: 0;
}
.captcha-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.captcha-preview {
  position: relative;
  width: 100%;
  height: 88px;
  border: 1px solid #dbe1ea;
  border-radius: 12px;
  background-color: #f8fafc;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s ease;
}
.captcha-preview:hover {
  border-color: #94a3b8;
  cursor: pointer;
}
.captcha-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.captcha-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #8a94a6;
}
.loading-icon {
  animation: spin 1s linear infinite;
  color: #6366f1;
}
.refresh-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 2px;
  font-size: 14px;
  color: #4c51bf;
}
.code-input {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  width: 100%;
}
.code-box {
  width: 42px;
  height: 42px;
  line-height: 42px;
  text-align: center;
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  background-color: var(--el-fill-color-blank);
  color: var(--el-text-color-regular);
  outline: none;
  font-size: 16px;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
.code-box:focus {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 1px var(--el-color-primary) inset;
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>


