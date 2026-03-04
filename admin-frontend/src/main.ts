import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import './styles/theme.css'
import './styles/global.scss'

import App from './App.vue'
import router from './router'
import { setupGlobalErrorHandler } from './composables/useErrorHandler'

const app = createApp(App)

// Pinia
const pinia = createPinia()
app.use(pinia)

// Router
app.use(router)

// Element Plus
app.use(ElementPlus)

// Element Plus Icons
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 初始化主题
const { updateTheme } = await import('./composables/useDarkMode')
updateTheme()

// 全局错误处理
setupGlobalErrorHandler(app, router)

app.mount('#app')
