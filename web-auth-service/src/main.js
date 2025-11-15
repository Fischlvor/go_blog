import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import Login from './views/Login.vue'
import Register from './views/Register.vue'
import ForgotPassword from './views/ForgotPassword.vue'
import QQCallback from './views/QQCallback.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/forgot-password', component: ForgotPassword },
  { path: '/oauth/qq/callback', component: QQCallback }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const pinia = createPinia()

createApp(App)
  .use(router)
  .use(pinia)
  .use(ElementPlus)
  .mount('#app')

