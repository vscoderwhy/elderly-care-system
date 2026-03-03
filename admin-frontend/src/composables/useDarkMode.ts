import { ref, watch, onMounted } from 'vue'

const THEME_KEY = 'app-theme'
export type ThemeMode = 'light' | 'dark' | 'auto'

const mode = ref<ThemeMode>(
  (localStorage.getItem(THEME_KEY) as ThemeMode) || 'auto'
)
const isDark = ref(false)

// 检测系统主题
const getSystemTheme = (): boolean => {
  if (typeof window === 'undefined') return false
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

const updateTheme = () => {
  const shouldDark = mode.value === 'dark' || (mode.value === 'auto' && getSystemTheme())
  isDark.value = shouldDark
  document.documentElement.setAttribute('data-theme', shouldDark ? 'dark' : 'light')

  // 加载 Element Plus 深色主题 CSS
  const loadDarkThemeCSS = () => {
    if (shouldDark) {
      // 动态加载 Element Plus 深色主题
      const linkId = 'element-dark-theme'
      let link = document.getElementById(linkId) as HTMLLinkElement
      if (!link) {
        link = document.createElement('link')
        link.id = linkId
        link.rel = 'stylesheet'
        link.href = 'https://unpkg.com/element-plus/theme-chalk/dark/css-vars.css'
        document.head.appendChild(link)
      }
    } else {
      // 移除深色主题 CSS
      const link = document.getElementById('element-dark-theme')
      if (link) {
        link.remove()
      }
    }
  }

  loadDarkThemeCSS()
}

export const useDarkMode = () => {
  // 初始化主题
  onMounted(() => {
    updateTheme()

    // 监听系统主题变化
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handleSystemThemeChange = () => {
      if (mode.value === 'auto') {
        updateTheme()
      }
    }
    mediaQuery.addEventListener('change', handleSystemThemeChange)

    // 返回清理函数
    return () => {
      mediaQuery.removeEventListener('change', handleSystemThemeChange)
    }
  })

  // 监听模式变化
  watch(mode, (newMode) => {
    localStorage.setItem(THEME_KEY, newMode)
    updateTheme()
  })

  const toggle = () => {
    if (mode.value === 'dark') {
      mode.value = 'light'
    } else if (mode.value === 'light') {
      mode.value = 'auto'
    } else {
      mode.value = 'dark'
    }
  }

  const setMode = (newMode: ThemeMode) => {
    mode.value = newMode
  }

  return {
    mode,
    isDark,
    toggle,
    setMode,
    updateTheme
  }
}
