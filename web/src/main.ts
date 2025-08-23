import {createApp} from 'vue'
import '@/assets/base.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import {pinia} from "@/stores";

// 全局导入md-editor-v3样式，避免重复导入导致的样式冲突
import 'md-editor-v3/lib/style.css'
import 'md-editor-v3/lib/preview.css'

const app = createApp(App)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.use(pinia).use(router)

app.mount('#app')