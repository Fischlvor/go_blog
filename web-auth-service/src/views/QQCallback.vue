<template>
  <div class="callback-container">
    <div class="callback-card">
      <p v-if="loading">正在处理QQ登录...</p>
      <p v-else class="success">登录成功，正在跳转...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { qqLogin } from '@/api/auth'
import { getDeviceId, getBrowserName } from '@/utils/device'

const route = useRoute()
const router = useRouter()

const loading = ref(true)

onMounted(async () => {
  // 从URL获取参数
  const code = route.query.code
  const appId = route.query.app_id
  const redirectUri = route.query.redirect_uri || ''
  const state = route.query.state || ''

  if (!code) {
    ElMessage.error('缺少授权码')
    loading.value = false
    return
  }

  if (!appId) {
    ElMessage.error('访问链接无效')
    loading.value = false
    return
  }

  try {
    const response = await qqLogin({
      code: code,
      app_id: appId,
      redirect_uri: redirectUri,
      device_id: getDeviceId(),
      device_name: getBrowserName(),
      device_type: 'web'
    })

    if (response.data.code === 0) {
      const data = response.data.data
      // 跳转到回调地址
      if (data.redirect_uri) {
        window.location.href = `${data.redirect_uri}?` +
          `code=${data.code}&` +
          `state=${state}`
      } else {
        ElMessage.error('缺少回调地址')
        loading.value = false
      }
    } else {
      ElMessage.error(response.data.message || 'QQ登录失败')
      loading.value = false
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || 'QQ登录失败')
    loading.value = false
  }
})

// getDeviceId 和 getBrowserName 已从 @/utils/device 导入
</script>

<style scoped>
.callback-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
}

.callback-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  text-align: center;
  min-width: 300px;
}

.success {
  color: #3c3;
}
</style>

