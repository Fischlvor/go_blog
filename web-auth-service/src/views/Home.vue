<template>
  <div class="home-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>
    
    <!-- 顶部导航 -->
    <div class="home-header">
      <div class="header-content">
        <div class="logo-section">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h1 class="logo-text">SSO 统一登录</h1>
        </div>
        <div class="header-actions">
          <el-button v-if="isLoggedIn" type="primary" @click="goToManage">管理中心</el-button>
          <el-button v-else type="primary" @click="goToLogin">登录</el-button>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="home-main">
      <!-- 左侧展示区 -->
      <div class="left-showcase">
        <div class="showcase-content">
          <div class="showcase-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h1 class="showcase-title">统一身份认证</h1>
          <p class="showcase-subtitle">一次登录，畅享所有应用</p>
          <div class="showcase-features">
            <div class="feature-item">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M9 12L11 14L15 10M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <span>安全可靠</span>
            </div>
            <div class="feature-item">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M13 10V3L4 14H11L11 21L20 10L13 10Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <span>快速便捷</span>
            </div>
            <div class="feature-item">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 15C13.6569 15 15 13.6569 15 12C15 10.3431 13.6569 9 12 9C10.3431 9 9 10.3431 9 12C9 13.6569 10.3431 15 12 15Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M19.4 15C19.2669 15.3016 19.2272 15.6362 19.286 15.9606C19.3448 16.285 19.4995 16.5843 19.73 16.82L19.79 16.88C19.976 17.0657 20.1235 17.2863 20.2241 17.5291C20.3248 17.7719 20.3766 18.0322 20.3766 18.295C20.3766 18.5578 20.3248 18.8181 20.2241 19.0609C20.1235 19.3037 19.976 19.5243 19.79 19.71C19.6043 19.896 19.3837 20.0435 19.1409 20.1441C18.8981 20.2448 18.6378 20.2966 18.375 20.2966C18.1122 20.2966 17.8519 20.2448 17.6091 20.1441C17.3663 20.0435 17.1457 19.896 16.96 19.71L16.9 19.65C16.6643 19.4195 16.365 19.2648 16.0406 19.206C15.7162 19.1472 15.3816 19.1869 15.08 19.32C14.7842 19.4468 14.532 19.6572 14.3543 19.9255C14.1766 20.1938 14.0813 20.5082 14.08 20.83V21C14.08 21.5304 13.8693 22.0391 13.4942 22.4142C13.1191 22.7893 12.6104 23 12.08 23C11.5496 23 11.0409 22.7893 10.6658 22.4142C10.2907 22.0391 10.08 21.5304 10.08 21V20.91C10.0723 20.579 9.96512 20.258 9.77251 19.9887C9.5799 19.7194 9.31074 19.5143 9 19.4C8.69838 19.2669 8.36381 19.2272 8.03941 19.286C7.71502 19.3448 7.41568 19.4995 7.18 19.73L7.12 19.79C6.93425 19.976 6.71368 20.1235 6.47088 20.2241C6.22808 20.3248 5.96783 20.3766 5.705 20.3766C5.44217 20.3766 5.18192 20.3248 4.93912 20.2241C4.69632 20.1235 4.47575 19.976 4.29 19.79C4.10405 19.6043 3.95653 19.3837 3.85588 19.1409C3.75523 18.8981 3.70343 18.6378 3.70343 18.375C3.70343 18.1122 3.75523 17.8519 3.85588 17.6091C3.95653 17.3663 4.10405 17.1457 4.29 16.96L4.35 16.9C4.58054 16.6643 4.73519 16.365 4.794 16.0406C4.85282 15.7162 4.81312 15.3816 4.68 15.08C4.55324 14.7842 4.34276 14.532 4.07447 14.3543C3.80618 14.1766 3.49179 14.0813 3.17 14.08H3C2.46957 14.08 1.96086 13.8693 1.58579 13.4942C1.21071 13.1191 1 12.6104 1 12.08C1 11.5496 1.21071 11.0409 1.58579 10.6658C1.96086 10.2907 2.46957 10.08 3 10.08H3.09C3.42099 10.0723 3.742 9.96512 4.0113 9.77251C4.28059 9.5799 4.48572 9.31074 4.6 9C4.73312 8.69838 4.77282 8.36381 4.714 8.03941C4.65519 7.71502 4.50054 7.41568 4.27 7.18L4.21 7.12C4.02405 6.93425 3.87653 6.71368 3.77588 6.47088C3.67523 6.22808 3.62343 5.96783 3.62343 5.705C3.62343 5.44217 3.67523 5.18192 3.77588 4.93912C3.87653 4.69632 4.02405 4.47575 4.21 4.29C4.39575 4.10405 4.61632 3.95653 4.85912 3.85588C5.10192 3.75523 5.36217 3.70343 5.625 3.70343C5.88783 3.70343 6.14808 3.75523 6.39088 3.85588C6.63368 3.95653 6.85425 4.10405 7.04 4.29L7.1 4.35C7.33568 4.58054 7.63502 4.73519 7.95941 4.794C8.28381 4.85282 8.61838 4.81312 8.92 4.68H9C9.29577 4.55324 9.54802 4.34276 9.72569 4.07447C9.90337 3.80618 9.99872 3.49179 10 3.17V3C10 2.46957 10.2107 1.96086 10.5858 1.58579C10.9609 1.21071 11.4696 1 12 1C12.5304 1 13.0391 1.21071 13.4142 1.58579C13.7893 1.96086 14 2.46957 14 3V3.09C14.0013 3.41179 14.0966 3.72618 14.2743 3.99447C14.452 4.26276 14.7042 4.47324 15 4.6C15.3016 4.73312 15.6362 4.77282 15.9606 4.714C16.285 4.65519 16.5843 4.50054 16.82 4.27L16.88 4.21C17.0657 4.02405 17.2863 3.87653 17.5291 3.77588C17.7719 3.67523 18.0322 3.62343 18.295 3.62343C18.5578 3.62343 18.8181 3.67523 19.0609 3.77588C19.3037 3.87653 19.5243 4.02405 19.71 4.21C19.896 4.39575 20.0435 4.61632 20.1441 4.85912C20.2448 5.10192 20.2966 5.36217 20.2966 5.625C20.2966 5.88783 20.2448 6.14808 20.1441 6.39088C20.0435 6.63368 19.896 6.85425 19.71 7.04L19.65 7.1C19.4195 7.33568 19.2648 7.63502 19.206 7.95941C19.1472 8.28381 19.1869 8.61838 19.32 8.92V9C19.4468 9.29577 19.6572 9.54802 19.9255 9.72569C20.1938 9.90337 20.5082 9.99872 20.83 10H21C21.5304 10 22.0391 10.2107 22.4142 10.5858C22.7893 10.9609 23 11.4696 23 12C23 12.5304 22.7893 13.0391 22.4142 13.4142C22.0391 13.7893 21.5304 14 21 14H20.91C20.5882 14.0013 20.2738 14.0966 20.0055 14.2743C19.7372 14.452 19.5268 14.7042 19.4 15Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <span>统一管理</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧应用区 -->
      <div class="right-apps">
        <div class="apps-header">
          <h2 class="apps-title">我的应用</h2>
          <p class="apps-subtitle">选择一个应用开始使用</p>
        </div>

        <!-- 应用列表 -->
        <div v-if="loading" class="loading-section">
          <el-icon class="is-loading"><Loading /></el-icon>
          <p>加载应用列表...</p>
        </div>

        <div v-else-if="applications.length === 0" class="empty-section">
          <el-empty description="暂无可用应用" />
        </div>

        <div v-else class="applications-grid">
          <div 
            v-for="app in applications" 
            :key="app.id" 
            class="app-card"
            @click="handleAppClick(app)"
          >
            <div class="app-icon">
              <img v-if="app.icon" :src="app.icon" :alt="app.app_name" />
              <div v-else class="default-icon">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <rect x="3" y="3" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/>
                  <path d="M9 9H15M9 12H15M9 15H12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
              </div>
            </div>
            <div class="app-name">{{ app.app_name }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <div class="home-footer">
      <p>&copy; 2025 SSO 统一登录系统</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { getPublicApplications } from '../api/oauth'

const router = useRouter()
const applications = ref([])
const loading = ref(true)
const isLoggedIn = ref(false)

// 检查是否已登录
const checkLoginStatus = () => {
  const token = localStorage.getItem('access_token')
  isLoggedIn.value = !!token
}

// 获取公开应用列表
const fetchApplications = async () => {
  try {
    loading.value = true
    const response = await getPublicApplications()
    
    if (response.data.code === 0) {
      applications.value = response.data.data || []
    } else {
      ElMessage.error(response.data.msg || '获取应用列表失败')
    }
  } catch (error) {
    console.error('获取应用列表失败:', error)
    ElMessage.error('获取应用列表失败')
  } finally {
    loading.value = false
  }
}

// 处理应用点击
const handleAppClick = (app) => {
  // redirect_uri 现在是后端根据环境返回的单个URL
  if (!app.redirect_uri) {
    ElMessage.error('应用未配置回调地址')
    return
  }
  
  // 获取 API 基础地址
  const baseURL = import.meta.env.VITE_BASE_API || 'http://localhost:8001'
  
  // 跳转到 OAuth 授权端点
  const authorizeUrl = `${baseURL}/oauth/authorize?app_id=${app.app_key}&redirect_uri=${encodeURIComponent(app.redirect_uri)}`
  
  window.location.href = authorizeUrl
}

// 跳转到管理中心
const goToManage = () => {
  router.push('/manage')
}

// 跳转到登录页面
const goToLogin = () => {
  router.push('/login?app_id=manage')
}

onMounted(() => {
  checkLoginStatus()
  fetchApplications()
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.home-header {
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  padding: 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1280px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo-icon {
  width: 32px;
  height: 32px;
  color: #4f46e5;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo-text {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #111827;
  letter-spacing: -0.02em;
}

.home-main {
  flex: 1;
  display: flex;
  width: 100%;
}

/* 左侧展示区 */
.left-showcase {
  flex: 0 0 45%;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4rem 3rem;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.04);
  position: relative;
  overflow: hidden;
}

.background-decoration {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 9999;
  pointer-events: none;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(66, 133, 244, 0.08) 0%, rgba(52, 168, 83, 0.08) 100%);
  animation: float 20s infinite ease-in-out;
}

.circle-1 {
  width: 400px;
  height: 400px;
  top: -150px;
  right: -150px;
  animation-delay: 0s;
}

.circle-2 {
  width: 300px;
  height: 300px;
  bottom: -100px;
  left: -100px;
  animation-delay: 7s;
}

.circle-3 {
  width: 200px;
  height: 200px;
  top: 40%;
  left: 15%;
  animation-delay: 14s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -30px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

.showcase-content {
  max-width: 500px;
  color: #202124;
  position: relative;
  z-index: 1;
}

.showcase-icon {
  width: 80px;
  height: 80px;
  margin-bottom: 2rem;
  color: #4285f4;
}

.showcase-icon svg {
  width: 100%;
  height: 100%;
}

.showcase-title {
  font-size: 3rem;
  font-weight: 700;
  margin: 0 0 1rem 0;
  letter-spacing: -0.03em;
  line-height: 1.2;
}

.showcase-subtitle {
  font-size: 1.25rem;
  margin: 0 0 3rem 0;
  color: #5f6368;
  line-height: 1.6;
}

.showcase-features {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 1.125rem;
  color: #202124;
}

.feature-item svg {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  color: #34a853;
}

/* 右侧应用区 */
.right-apps {
  flex: 1;
  padding: 3rem 4rem;
  overflow-y: auto;
  background: white;
}

.apps-header {
  margin-bottom: 2rem;
}

.apps-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #111827;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.apps-subtitle {
  font-size: 1rem;
  color: #6b7280;
  margin: 0;
}

.loading-section {
  text-align: center;
  padding: 4rem 0;
  color: #6b7280;
}

.loading-section .el-icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  color: #4f46e5;
}

.loading-section p {
  font-size: 1rem;
}

.empty-section {
  background: white;
  border-radius: 16px;
  padding: 4rem 2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

.applications-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 1.5rem;
  margin-top: 2rem;
}

.app-card {
  background: white;
  border-radius: 16px;
  padding: 1.75rem 1.5rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.app-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(79, 70, 229, 0.15);
}

.app-icon {
  width: 100%;
  height: 80px;
  margin: 0 auto 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-icon img {
  max-width: 120px;
  max-height: 64px;
  object-fit: contain;
}

.default-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f1f3f4;
  border-radius: 14px;
  color: #5f6368;
}

.default-icon svg {
  width: 60%;
  height: 60%;
}

@media (max-width: 768px) {
  .default-icon {
    width: 56px;
    height: 56px;
  }
}

.app-name {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #111827;
  margin: 0;
  line-height: 1.4;
}

.home-footer {
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.04);
  padding: 1.5rem 0;
  text-align: center;
  color: #9ca3af;
  font-size: 0.875rem;
}

.home-footer p {
  margin: 0;
}

@media (max-width: 1024px) {
  .home-main {
    flex-direction: column;
  }
  
  .left-showcase {
    flex: 0 0 auto;
    padding: 3rem 2rem;
  }
  
  .showcase-title {
    font-size: 2.5rem;
  }
  
  .showcase-subtitle {
    font-size: 1.125rem;
  }
  
  .right-apps {
    padding: 2rem;
  }
}

@media (max-width: 768px) {
  .header-content {
    padding: 0.875rem 1rem;
  }
  
  .logo-text {
    font-size: 1.125rem;
  }
  
  .left-showcase {
    padding: 2rem 1.5rem;
  }
  
  .showcase-icon {
    width: 60px;
    height: 60px;
    margin-bottom: 1.5rem;
  }
  
  .showcase-title {
    font-size: 2rem;
  }
  
  .showcase-subtitle {
    font-size: 1rem;
    margin-bottom: 2rem;
  }
  
  .feature-item {
    font-size: 1rem;
  }
  
  .right-apps {
    padding: 1.5rem 1rem;
  }
  
  .apps-title {
    font-size: 1.5rem;
  }
  
  .apps-subtitle {
    font-size: 0.9375rem;
  }
  
  .applications-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 1rem;
  }
  
  .app-card {
    padding: 1.5rem 1.25rem;
  }
  
  .app-icon {
    height: 56px;
  }
  
  .app-name {
    font-size: 0.875rem;
  }
  
  .default-icon {
    width: 56px;
    height: 56px;
  }
}

@media (min-width: 1024px) {
  .applications-grid {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 2rem;
  }
}
</style>
