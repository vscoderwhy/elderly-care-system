<template>
  <!-- 全屏页面（如指挥中心） -->
  <router-view v-if="route.meta.fullScreen"></router-view>

  <!-- 正常布局 -->
  <div v-else class="layout-container" :class="{ 'is-mobile': isMobile }">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapsed ? '64px' : '200px'" class="layout-aside">
      <div class="logo-container" :class="{ 'is-collapsed': isCollapsed }">
        <img src="/logo.svg" alt="Logo" class="logo-icon" />
        <transition name="logo-fade">
          <span v-show="!isCollapsed" class="logo-title">养老院管理系统</span>
        </transition>
      </div>

      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapsed"
        :unique-opened="true"
        router
        class="layout-menu"
      >
        <el-menu-item index="/dashboard/overview">
          <el-icon><DataBoard /></el-icon>
          <template #title>概览仪表盘</template>
        </el-menu-item>

        <el-menu-item index="/dashboard/command">
          <el-icon><Monitor /></el-icon>
          <template #title>指挥中心大屏</template>
        </el-menu-item>

        <el-sub-menu index="statistics">
          <template #title>
            <el-icon><DataAnalysis /></el-icon>
            <span>数据分析</span>
          </template>
          <el-menu-item index="/statistics/advanced">高级看板</el-menu-item>
          <el-menu-item index="/statistics/report">报表中心</el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/monitor/realtime">
          <el-icon><VideoCamera /></el-icon>
          <template #title>实时监控</template>
        </el-menu-item>

        <el-sub-menu index="elderly">
          <template #title>
            <el-icon><User /></el-icon>
            <span>老人管理</span>
          </template>
          <el-menu-item index="/elderly/list">老人列表</el-menu-item>
          <el-menu-item index="/elderly/health">健康档案</el-menu-item>
          <el-menu-item index="/elderly/care">护理记录</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="nursing">
          <template #title>
            <el-icon><Briefcase /></el-icon>
            <span>护理管理</span>
          </template>
          <el-menu-item index="/nursing/tasks">护理任务</el-menu-item>
          <el-menu-item index="/nursing/records">护理记录</el-menu-item>
          <el-menu-item index="/nursing/assessment">护理评估</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="medical">
          <template #title>
            <el-icon><FirstAidKit /></el-icon>
            <span>医务管理</span>
          </template>
          <el-menu-item index="/medical/prescription">处方管理</el-menu-item>
          <el-menu-item index="/medical/rounds">查房记录</el-menu-item>
          <el-menu-item index="/medical/emergency">急救记录</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="finance">
          <template #title>
            <el-icon><Wallet /></el-icon>
            <span>财务管理</span>
          </template>
          <el-menu-item index="/finance/bills">费用账单</el-menu-item>
          <el-menu-item index="/finance/payments">支付记录</el-menu-item>
          <el-menu-item index="/finance/refund">退款管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="hr">
          <template #title>
            <el-icon><UserFilled /></el-icon>
            <span>人事管理</span>
          </template>
          <el-menu-item index="/hr/staff">员工管理</el-menu-item>
          <el-menu-item index="/hr/schedule">排班管理</el-menu-item>
          <el-menu-item index="/hr/attendance">考勤管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="inventory">
          <template #title>
            <el-icon><Box /></el-icon>
            <span>库存管理</span>
          </template>
          <el-menu-item index="/inventory/supplies">用品管理</el-menu-item>
          <el-menu-item index="/inventory/equipment">设备管理</el-menu-item>
          <el-menu-item index="/inventory/purchase">采购管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="visits">
          <template #title>
            <el-icon><Phone /></el-icon>
            <span>访视管理</span>
          </template>
          <el-menu-item index="/visits/appointment">预约管理</el-menu-item>
          <el-menu-item index="/visits/records">访视记录</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </template>
          <el-menu-item index="/system/users">用户管理</el-menu-item>
          <el-menu-item index="/system/roles">角色权限</el-menu-item>
          <el-menu-item index="/system/logs">操作日志</el-menu-item>
          <el-menu-item index="/system/settings">系统配置</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container class="layout-main">
      <!-- 头部 -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-button
            :icon="isCollapsed ? Expand : Fold"
            circle
            @click="toggleCollapse"
            class="collapse-btn"
          />
          <el-breadcrumb separator="/">
            <el-breadcrumb-item
              v-for="item in breadcrumbs"
              :key="item.path"
              :to="item.path"
            >
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="header-right">
          <!-- 搜索框 -->
          <el-input
            v-model="searchKeyword"
            placeholder="搜索..."
            prefix-icon="Search"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          />

          <!-- 主题切换器 -->
          <el-dropdown trigger="click" @command="handleThemeChange">
            <el-button circle>
              <el-icon>
                <Sunny v-if="darkMode.mode === 'light'" />
                <Moon v-else-if="darkMode.mode === 'dark'" />
                <Monitor v-else />
              </el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="light">
                  <el-icon><Sunny /></el-icon>
                  浅色模式
                </el-dropdown-item>
                <el-dropdown-item command="dark">
                  <el-icon><Moon /></el-icon>
                  深色模式
                </el-dropdown-item>
                <el-dropdown-item command="auto">
                  <el-icon><Monitor /></el-icon>
                  跟随系统
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <!-- 全屏按钮 -->
          <el-button :icon="isFullscreen ? ExitFullScreen : FullScreen" circle @click="toggleFullscreen" />

          <!-- 通知 -->
          <el-badge :value="unreadCount" :hidden="unreadCount === 0">
            <el-button circle>
              <el-icon><Bell /></el-icon>
            </el-button>
          </el-badge>

          <!-- 用户菜单 -->
          <el-dropdown trigger="click">
            <div class="user-avatar-wrapper">
              <el-avatar :size="32" :src="userAvatar">
                {{ userName.charAt(0) }}
              </el-avatar>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item>
                  <el-icon><Setting /></el-icon>
                  设置
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main class="layout-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>

      <!-- 页脚 -->
      <el-footer v-if="showFooter" class="layout-footer">
        <div class="footer-content">
          <span>© 2026 养老院管理系统 v1.0.0</span>
          <span>
            <el-link type="primary" href="https://github.com" target="_blank">
              GitHub
            </el-link>
            <el-divider direction="vertical" />
            <el-link type="primary" @click="handleAbout">关于</el-link>
          </span>
        </div>
      </el-footer>
    </el-container>

    <!-- 移动端抽屉菜单 -->
    <el-drawer
      v-model="mobileMenuVisible"
      direction="ltr"
      :with-header="false"
      size="200px"
      class="mobile-drawer"
    >
      <div class="mobile-menu-content">
        <!-- 侧边栏内容 -->
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DataBoard,
  DataAnalysis,
  VideoCamera,
  Monitor,
  User,
  Briefcase,
  FirstAidKit,
  Wallet,
  UserFilled,
  Box,
  Phone,
  Setting,
  Fold,
  Expand,
  Search,
  Sunny,
  Moon,
  Monitor,
  FullScreen,
  ExitFullScreen,
  Bell,
  SwitchButton
} from '@element-plus/icons-vue'
import { useDarkMode } from '@/composables/useDarkMode'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()

// 侧边栏状态
const isCollapsed = ref(false)
const mobileMenuVisible = ref(false)

// 响应式
const isMobile = ref(false)
const showFooter = ref(true)

// 用户信息
const userName = ref('管理员')
const userAvatar = ref('')
const unreadCount = ref(5)

// 搜索
const searchKeyword = ref('')

// 暗黑模式
const darkMode = useDarkMode()

// 全屏状态
const isFullscreen = ref(false)

// 当前激活的菜单
const activeMenu = computed(() => route.path)

// 面包屑
const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  return matched.map(item => ({
    path: item.path,
    title: item.meta?.title || ''
  }))
})

// 切换侧边栏折叠
const toggleCollapse = () => {
  if (isMobile.value) {
    mobileMenuVisible.value = !mobileMenuVisible.value
  } else {
    isCollapsed.value = !isCollapsed.value
  }
}

// 主题切换
const handleThemeChange = (command: 'light' | 'dark' | 'auto') => {
  darkMode.setMode(command)
}

// 全屏切换
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

// 搜索
const handleSearch = () => {
  console.log('搜索', searchKeyword.value)
}

// 退出登录
const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    console.log('退出登录')
    router.push('/login')
  })
}

// 关于
const handleAbout = () => {
  console.log('关于')
}

// 检测是否为移动设备
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
  if (isMobile.value) {
    isCollapsed.value = true
  }
}

// 监听窗口大小变化
const handleResize = () => {
  checkMobile()
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', handleResize)

  // 监听全屏变化
  document.addEventListener('fullscreenchange', () => {
    isFullscreen.value = !!document.fullscreenElement
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

// 监听路由变化，更新面包屑
watch(
  () => route.path,
  () => {
    if (isMobile.value) {
      mobileMenuVisible.value = false
    }
  }
)
</script>

<style scoped lang="scss">
.layout-container {
  display: flex;
  height: 100vh;
  overflow: hidden;

  &.is-mobile {
    .layout-aside {
      position: fixed;
      z-index: 1000;
      height: 100vh;
    }
  }
}

.layout-aside {
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  transition: width 0.3s ease;
  overflow-x: hidden;
  overflow-y: auto;

  .logo-container {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    height: 60px;
    box-sizing: border-box;
    border-bottom: 1px solid var(--border-color);

    &.is-collapsed {
      justify-content: center;
      padding: 16px 0;
    }

    .logo-icon {
      width: 32px;
      height: 32px;
      flex-shrink: 0;
    }

    .logo-title {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      white-space: nowrap;
    }
  }

  .layout-menu {
    border-right: none;
  }
}

.layout-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.layout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

  .header-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .collapse-btn {
    flex-shrink: 0;
  }

  .search-input {
    width: 200px;
  }

  .user-avatar-wrapper {
    cursor: pointer;
    transition: var(--transition-base);

    &:hover {
      opacity: 0.8;
    }
  }
}

.layout-content {
  flex: 1;
  overflow-y: auto;
  background: var(--bg-secondary);
}

.layout-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background: var(--card-bg);
  border-top: 1px solid var(--border-color);
  height: 50px;

  .footer-content {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 12px;
    color: var(--text-secondary);
  }
}

.mobile-drawer {
  :deep(.el-drawer__body) {
    padding: 0;
  }
}

.mobile-menu-content {
  height: 100%;
}

// 动画
.logo-fade-enter-active,
.logo-fade-leave-active {
  transition: opacity 0.3s;
}

.logo-fade-enter-from,
.logo-fade-leave-to {
  opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateX(10px);
}

// 响应式适配
@media (max-width: 768px) {
  .layout-header {
    padding: 0 12px;

    .header-left {
      gap: 8px;
    }

    .header-right {
      gap: 8px;
    }

    .search-input {
      width: 120px;
    }

    :deep(.el-breadcrumb) {
      display: none;
    }
  }

  .layout-footer {
    padding: 0 12px;

    .footer-content {
      flex-direction: column;
      gap: 4px;
      text-align: center;
    }
  }

  .layout-aside {
    width: 0 !important;

    &.el-aside {
      width: 200px !important;
    }
  }
}

@media (max-width: 480px) {
  .layout-header {
    .search-input {
      display: none;
    }
  }
}
</style>
