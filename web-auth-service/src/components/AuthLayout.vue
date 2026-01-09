<template>
  <div class="auth-container" :data-app-theme="appTheme">
    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>
    <div class="auth-card">
      <div class="logo-section">
        <div class="logo-icon">
          <slot name="icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </slot>
        </div>
        <div class="logo-text">
          <h1 class="title">{{ title }}</h1>
          <p class="subtitle">{{ subtitle }}</p>
        </div>
      </div>
      <slot></slot>
      <div class="security-tip">
        <svg class="tip-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 8V12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 16H12.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>{{ securityTip }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
  securityTip: { type: String, default: '您的信息已加密传输，请放心使用' },
  appId: { type: String, default: 'blog' }
})

const appTheme = computed(() => {
  const themes = { 'blog': 'blue', 'mcp': 'green', 'manage': 'purple' }
  return themes[props.appId] || 'purple'
})
</script>

<style scoped>
.auth-container {
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
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
}

.auth-card {
  position: relative;
  z-index: 1;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(255, 255, 255, 0.5);
  padding: 32px;
  width: 100%;
  max-width: min(480px, calc(100vw - 40px));
  min-width: 0;
  box-sizing: border-box;
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

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
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.logo-icon svg,
.logo-icon :deep(svg) {
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

.security-tip {
  margin-top: 24px;
  padding: 12px 16px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 12px;
  color: #6b7280;
}

.tip-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  color: #667eea;
}

/* 表单公共样式 */
:deep(.auth-form),
:deep(.register-form),
:deep(.login-form),
:deep(.forgot-password-form) {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

:deep(.form-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

:deep(.form-group label) {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

:deep(.input-icon) {
  width: 18px;
  height: 18px;
  color: #9ca3af;
  flex-shrink: 0;
}

:deep(.form-input) {
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

:deep(.form-input::placeholder) {
  color: #9ca3af;
}

:deep(.form-input:focus) {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

:deep(.form-input:disabled) {
  background-color: #f3f4f6;
  color: #9ca3af;
  cursor: not-allowed;
  opacity: 0.6;
}

:deep(.error-text) {
  color: #ef4444;
  font-size: 12px;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

:deep(.verification-group) {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  width: 100%;
  min-width: 0;
}

:deep(.verification-group .form-input) {
  flex: 1;
}

:deep(.captcha-group) {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  width: 100%;
  min-width: 0;
}

:deep(.captcha-group .form-input) {
  flex: 1;
}

:deep(.captcha-wrapper) {
  width: 130px;
  height: 48px;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid #e5e7eb;
  transition: all 0.3s;
  flex-shrink: 0;
}

:deep(.captcha-wrapper:hover) {
  border-color: #667eea;
}

:deep(.captcha-image) {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

:deep(.captcha-placeholder) {
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

:deep(.captcha-placeholder span) {
  color: #999;
  font-size: 12px;
}

:deep(.btn-primary),
:deep(.btn-register),
:deep(.btn-login),
:deep(.btn-submit) {
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

:deep(.btn-primary:hover:not(:disabled)),
:deep(.btn-register:hover:not(:disabled)),
:deep(.btn-login:hover:not(:disabled)),
:deep(.btn-submit:hover:not(:disabled)) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.5);
}

:deep(.btn-primary:disabled),
:deep(.btn-register:disabled),
:deep(.btn-login:disabled),
:deep(.btn-submit:disabled) {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
  background: #d1d5db !important;
  color: #9ca3af !important;
}

:deep(.btn-send-code) {
  padding: 0 20px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
  flex-shrink: 0;
}

:deep(.btn-send-code:hover:not(:disabled)) {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(102, 126, 234, 0.4);
}

:deep(.btn-send-code:disabled) {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
  background: #d1d5db !important;
  color: #9ca3af !important;
}

:deep(.links) {
  text-align: center;
}

:deep(.link-item) {
  color: #667eea;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
}

:deep(.link-item:hover) {
  color: #764ba2;
  text-decoration: underline;
}

:deep(.divider-dot) {
  color: #d1d5db;
  margin: 0 8px;
}

:deep(.loading-content) {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

:deep(.loading-dots) {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
