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
const logout = () => {
  localStorage.removeItem('access_token')
  router.push('/login')
  ElMessage.success('已退出登录')
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
