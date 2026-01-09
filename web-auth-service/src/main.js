import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import Home from './views/Home.vue'
import Login from './views/Login.vue'
import Register from './views/Register.vue'
import ForgotPassword from './views/ForgotPassword.vue'
import QQCallback from './views/QQCallback.vue'
import ManageLayout from './views/manage/Layout.vue'
import Devices from './views/manage/Devices.vue'
import Security from './views/manage/Security.vue'
import Activity from './views/manage/Activity.vue'

const routes = [
  { path: '/', component: Home, name: 'Home' },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/forgot-password', component: ForgotPassword },
  { path: '/oauth/qq/callback', component: QQCallback },
  
  // 管理后台路由
  {
    path: '/manage',
    component: ManageLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/manage/devices' },
      { path: 'devices', component: Devices, name: 'DeviceManage' },
      { path: 'security', component: Security, name: 'SecuritySettings' },
      { path: 'activity', component: Activity, name: 'ActivityLogs' }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 检查 /login 路由的 app_id 参数
  if (to.path === '/login') {
    const appId = to.query.app_id
    const validAppIds = ['blog', 'mcp', 'manage']
    
    // 没有 app_id 或 app_id 不在有效列表中，重定向到主页
    if (!appId || !validAppIds.includes(appId)) {
      next('/')
      return
    }
  }
  
  if (to.meta.requiresAuth) {
    // 检查是否已登录
    const token = localStorage.getItem('access_token')
    if (!token) {
      next('/login')
      return
    }
  }
  next()
})

const pinia = createPinia()

createApp(App)
  .use(router)
  .use(pinia)
  .use(ElementPlus)
  .mount('#app')

