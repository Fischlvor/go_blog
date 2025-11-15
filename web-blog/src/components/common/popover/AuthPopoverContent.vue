<template>
  <div class="user-popover-content" @click.stop>
    <div class="title">
      <el-row>请登录以获取完整的功能体验。</el-row>
    </div>
    <div class="auth-button">
      <div class="button-item">
        <el-button class="login" icon="User" @click.stop="redirectToSSO('login')"/>
        登录
      </div>

      <div class="button-item">
        <el-button class="register" icon="Edit" @click.stop="redirectToSSO('register')"/>
        注册
      </div>

      <div class="button-item">
        <el-button class="forgot-password" icon="Unlock" @click.stop="redirectToSSO('forgot-password')"/>
        找回密码
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLayoutStore } from "@/stores/layout";
import axios from "axios";
import { ElMessage } from "element-plus";

const layoutStore = useLayoutStore();

// 跳转到SSO登录/注册/找回密码
const redirectToSSO = async (action: 'login' | 'register' | 'forgot-password') => {
  try {
    // 关闭弹窗
    layoutStore.hide('popoverVisible');
    
    // 构建回调地址
    const redirectUri = encodeURIComponent(window.location.origin + '/sso-callback');
    
    // 获取SSO登录URL
    const response = await axios.get(`/api/auth/sso_login_url?redirect_uri=${redirectUri}`);
    
    if (response.data.code === 0) {
      // 构建完整的SSO URL，根据action添加不同的路径
      let ssoURL = response.data.data.sso_login_url;
      
      // 根据action修改URL路径
      if (action === 'register') {
        ssoURL = ssoURL.replace('/login', '/register');
      } else if (action === 'forgot-password') {
        ssoURL = ssoURL.replace('/login', '/forgot-password');
      }
      
      // 跳转到SSO页面
      window.location.href = ssoURL;
    } else {
      ElMessage.error(response.data.message || '获取登录地址失败');
    }
  } catch (error: any) {
    console.error('获取SSO登录URL失败:', error);
    ElMessage.error(error.response?.data?.message || '登录服务异常，请稍后重试');
  }
};
</script>

<style scoped lang="scss">
.user-popover-content {

  .title {
    display: flex;
    height: 60px;

    .el-row {
      margin: auto auto;
    }
  }

  .auth-button {
    display: flex;
    justify-content: center;

    .button-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      width: 80px;

      .el-button {
        border: none;
        --el-font-size-base: 32px;
        height: 62px;
        background-color: transparent;
      }
    }
  }

}
</style>