<template>
  <div class="callback-container">
    <div class="callback-box">
      <el-icon class="is-loading" color="#667eea" :size="48">
        <Loading />
      </el-icon>
      <div class="callback-text">正在登录中...</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { Loading } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

onMounted(async () => {
  try {
    // 获取URL参数
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const state = urlParams.get('state')

    if (!code) {
      ElMessage.error('登录失败：未获取到授权码')
      router.push({ name: 'index' })
      return
    }

    // ✅ OAuth 2.0: 用code向应用后端换取token
    const response = await fetch(`/api/auth/callback?code=${code}&state=${state}&redirect_uri=${encodeURIComponent(window.location.origin + '/sso-callback')}`)
    const data = await response.json()

    if (data.code !== 0) {
      ElMessage.error('登录失败：' + data.message)
      router.push({ name: 'index' })
      return
    }

    // 保存access_token到store（refresh_token在后端session）
    userStore.state.accessToken = data.data.access_token
    userStore.state.isUserLoggedInBefore = true

    // ✅ 重置标志，确保能重新获取用户信息
    userStore.state.userInfoInitialized = false
    
    // 获取用户信息
    await userStore.initializeUserInfo()

    ElMessage.success('登录成功')
    
    // ✅ 使用state从sessionStorage读取返回路径
    let returnUrl = '/'
    if (state) {
      const stateKey = `oauth_state_${state}`
      const stateData = sessionStorage.getItem(stateKey)
      
      if (stateData) {
        try {
          const parsed = JSON.parse(stateData)
          returnUrl = parsed.returnUrl || '/'
          sessionStorage.removeItem(stateKey) // 清除已使用的state数据
        } catch (e) {
          console.error('解析state数据失败:', e)
        }
      }
    }
    
    // 跳转到返回页面或首页
    router.push(returnUrl)
  } catch (error) {
    console.error('SSO回调处理失败:', error)
    ElMessage.error('登录处理失败，请重试')
    router.push({ name: 'index' })
  }
})
</script>

<style scoped lang="scss">
.callback-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.callback-box {
  text-align: center;
  padding: 60px;
  background: white;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.callback-text {
  margin-top: 24px;
  font-size: 18px;
  color: #666;
  font-weight: 500;
}
</style>

