<template>
  <el-container class="manage-container">
    <!-- 顶部导航 -->
    <el-header class="manage-header">
      <div class="header-left">
        <h2>SSO管理中心</h2>
      </div>
      <div class="header-right">
        <el-dropdown @command="handleCommand">
          <span class="user-info">
            <el-avatar :size="32" />
            {{ userInfo.nickname || userInfo.email }}
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    
    <el-container>
      <!-- 侧边导航 -->
      <el-aside width="200px" class="manage-aside">
        <el-menu 
          :default-active="activeMenu" 
          router
          class="manage-menu"
        >
          <el-menu-item index="/manage/devices">
            <el-icon><Monitor /></el-icon>
            <span>设备管理</span>
          </el-menu-item>
          <el-menu-item index="/manage/security">
            <el-icon><Lock /></el-icon>
            <span>安全设置</span>
          </el-menu-item>
          <el-menu-item index="/manage/activity">
            <el-icon><Document /></el-icon>
            <span>操作日志</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <!-- 主内容区 -->
      <el-main class="manage-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Monitor, Lock, Document, ArrowDown } from '@element-plus/icons-vue'
import { manageApi } from '../../api/manage'

const router = useRouter()
const route = useRoute()

const userInfo = ref({})

const activeMenu = computed(() => route.path)

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const res = await manageApi.getProfile()
    if (res.code === 0) {
      userInfo.value = res.data
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

// 处理下拉菜单命令
const handleCommand = (command) => {
  if (command === 'logout') {
    logout()
  }
}

// 退出登录
const logout = async () => {
  try {
    const res = await manageApi.logout()
    console.log('manage应用退出响应:', res) // 调试信息
    if (res.data && res.data.code === 0) {
      ElMessage.success('退出成功')
      // 清除所有本地token
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      
      // 检查是否有重定向信息
      const redirectTo = res.data.data?.redirect_to
      if (redirectTo) {
        router.push(redirectTo) // 使用后端返回的重定向路径
      } else {
        router.push('/login?app_id=manage') // 默认跳转
      }
    } else {
      console.error('响应格式错误:', res) // 调试信息
      ElMessage.error(res.data?.message || 'SSO退出失败')
    }
  } catch (error) {
    console.error('SSO退出失败:', error)
    ElMessage.error('SSO退出失败')
    // 即使API失败，也清除所有本地token并跳转
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    router.push('/login?app_id=manage')
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.manage-container {
  height: 100vh;
}

.manage-header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.header-left h2 {
  margin: 0;
  color: #303133;
  font-size: 20px;
  font-weight: 500;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
}

.user-info:hover {
  color: #409eff;
}

.manage-aside {
  background: #f5f7fa;
  border-right: 1px solid #e4e7ed;
}

.manage-menu {
  border: none;
  background: transparent;
}

.manage-main {
  background: #f5f7fa;
  padding: 20px;
}
</style>
