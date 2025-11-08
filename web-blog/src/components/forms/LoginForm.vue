<template>
  <div class="login-dialog-container">
    <div class="login-form">
      <div class="sso-login-box">
        <div class="sso-title">统一登录</div>
        <div class="sso-description">为了更好的用户体验，我们采用统一登录系统</div>
        <el-button
            type="primary"
            size="large"
            class="sso-button"
            @click="redirectToSSO"
        >
          前往登录
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios';
import { useLayoutStore } from "@/stores/layout";

const layoutStore = useLayoutStore()

// 跳转到SSO登录
const redirectToSSO = async () => {
  try {
    // 获取SSO登录URL
    const redirectUri = encodeURIComponent(window.location.origin + '/sso-callback')
    const response = await axios.get(`/api/auth/sso_login_url?redirect_uri=${redirectUri}`)
    
    if (response.data.code === 0) {
      // 关闭登录弹窗
      layoutStore.state.loginVisible = false
      // 跳转到SSO登录页
      window.location.href = response.data.data.sso_login_url
    } else {
      ElMessage.error('获取登录地址失败')
    }
  } catch (error) {
    console.error('获取SSO登录URL失败:', error)
    ElMessage.error('登录服务异常，请稍后重试')
  }
}
</script>


<style scoped lang="scss">
.login-form {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.sso-login-box {
  text-align: center;
  padding: 40px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  min-width: 360px;
}

.sso-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.sso-description {
  font-size: 14px;
  color: #666;
  margin-bottom: 32px;
  line-height: 1.6;
}

.sso-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  }

  &:active {
    transform: translateY(0);
  }
}
</style>


