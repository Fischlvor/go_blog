<template>
  <div class="security-page">
    <div class="page-header">
      <h3>安全设置</h3>
    </div>

    <!-- 安全操作 -->
    <el-card class="security-card">
      <template #header>
        <div class="card-header">
          <el-icon><Lock /></el-icon>
          <span>安全操作</span>
        </div>
      </template>
      
      <div class="security-actions">
        <div class="action-item">
          <div class="action-info">
            <h4>退出当前设备SSO</h4>
            <p>退出当前设备的SSO登录状态，需要重新登录才能访问其他应用</p>
          </div>
          <el-button 
            type="warning" 
            @click="confirmSSOLogout"
            :loading="ssoLogouting"
          >
            退出SSO
          </el-button>
        </div>

        <el-divider />

        <div class="action-item">
          <div class="action-info">
            <h4>退出所有设备</h4>
            <p>强制退出所有设备上的登录状态，包括当前设备</p>
          </div>
          <el-button 
            type="danger" 
            @click="confirmLogoutAll"
            :loading="logouttingAll"
          >
            退出所有设备
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 安全提示 -->
    <el-card class="tips-card">
      <template #header>
        <div class="card-header">
          <el-icon><InfoFilled /></el-icon>
          <span>安全提示</span>
        </div>
      </template>
      
      <div class="security-tips">
        <el-alert
          title="关于SSO退出"
          type="info"
          :closable="false"
          show-icon
        >
          <template #default>
            <ul>
              <li>退出当前设备SSO：只影响当前设备，其他设备仍可正常使用</li>
              <li>退出所有设备：所有设备都需要重新登录</li>
              <li>如果怀疑账户被盗用，建议选择"退出所有设备"</li>
            </ul>
          </template>
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Lock, InfoFilled } from '@element-plus/icons-vue'
import { manageApi } from '../../api/manage'

const router = useRouter()

const ssoLogouting = ref(false)
const logouttingAll = ref(false)

// 确认SSO退出
const confirmSSOLogout = () => {
  ElMessageBox.confirm(
    '退出SSO后，您需要重新登录才能访问其他应用。是否继续？',
    '确认退出SSO',
    {
      confirmButtonText: '确认退出',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    await ssoLogout()
  }).catch(() => {
    // 用户取消
  })
}

// SSO退出
const ssoLogout = async () => {
  ssoLogouting.value = true
  try {
    const res = await manageApi.ssoLogout()
    console.log('SSO退出响应:', res) // 调试信息
    if (res.data && res.data.code === 0) {
      ElMessage.success('SSO退出成功')
      // 清除所有本地token并跳转到管理后台登录页
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      router.push('/login?app_id=manage')
    } else {
      console.error('响应格式错误:', res) // 调试信息
      ElMessage.error(res.data?.message || 'SSO退出失败')
    }
  } catch (error) {
    console.error('SSO退出失败:', error)
    ElMessage.error('SSO退出失败')
  } finally {
    ssoLogouting.value = false
  }
}

// 确认退出所有设备
const confirmLogoutAll = () => {
  ElMessageBox.confirm(
    '此操作将强制退出所有设备上的登录状态，包括当前设备。是否继续？',
    '确认退出所有设备',
    {
      confirmButtonText: '确认退出',
      cancelButtonText: '取消',
      type: 'error',
    }
  ).then(async () => {
    await logoutAll()
  }).catch(() => {
    // 用户取消
  })
}

// 退出所有设备
const logoutAll = async () => {
  logouttingAll.value = true
  try {
    const res = await manageApi.logoutAll()
    console.log('退出所有设备响应:', res) // 调试信息
    if (res.data && res.data.code === 0) {
      const kickedCount = res.data.data?.kicked_count || 0
      ElMessage.success(`已退出 ${kickedCount} 台设备`)
      // 清除本地token并跳转到管理后台登录页
      localStorage.removeItem('access_token')
      router.push('/login?app_id=manage')
    } else {
      console.error('响应格式错误:', res) // 调试信息
      ElMessage.error(res.data?.message || '退出所有设备失败')
    }
  } catch (error) {
    console.error('退出所有设备失败:', error)
    ElMessage.error('退出所有设备失败')
  } finally {
    logouttingAll.value = false
  }
}
</script>

<style scoped>
.security-page {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h3 {
  margin: 0;
  color: #303133;
  font-size: 18px;
  font-weight: 500;
}

.security-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #303133;
}

.security-actions {
  padding: 10px 0;
}

.action-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
}

.action-info h4 {
  margin: 0 0 8px 0;
  color: #303133;
  font-size: 16px;
  font-weight: 500;
}

.action-info p {
  margin: 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
}

.tips-card {
  margin-bottom: 20px;
}

.security-tips :deep(.el-alert__content) {
  padding-left: 0;
}

.security-tips ul {
  margin: 8px 0 0 0;
  padding-left: 20px;
}

.security-tips li {
  margin-bottom: 4px;
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
}
</style>
