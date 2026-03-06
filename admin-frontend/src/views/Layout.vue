<template>
  <el-container class="layout-container">
    <el-aside width="200px">
      <div class="logo">养老院管理</div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataBoard /></el-icon>
          <span>工作台</span>
        </el-menu-item>
        <el-menu-item index="/statistics">
          <el-icon><TrendCharts /></el-icon>
          <span>数据报表</span>
        </el-menu-item>
        <el-menu-item index="/elderly">
          <el-icon><User /></el-icon>
          <span>老人管理</span>
        </el-menu-item>
        <el-menu-item index="/care">
          <el-icon><Notebook /></el-icon>
          <span>护理记录</span>
        </el-menu-item>
        <el-menu-item index="/medications">
          <el-icon><Suitcase /></el-icon>
          <span>用药管理</span>
        </el-menu-item>
        <el-menu-item index="/visits">
          <el-icon><Calendar /></el-icon>
          <span>探视预约</span>
        </el-menu-item>
        <el-menu-item index="/alerts">
          <el-icon><Warning /></el-icon>
          <span>智能预警</span>
        </el-menu-item>
        <el-menu-item index="/bills">
          <el-icon><Wallet /></el-icon>
          <span>财务管理</span>
        </el-menu-item>
        <el-menu-item index="/export">
          <el-icon><Download /></el-icon>
          <span>数据导出</span>
        </el-menu-item>
        <el-menu-item index="/staff">
          <el-icon><Avatar /></el-icon>
          <span>员工管理</span>
        </el-menu-item>
        <el-menu-item index="/rooms">
          <el-icon><House /></el-icon>
          <span>房间管理</span>
        </el-menu-item>
        <el-sub-menu index="system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/system/users">
            <el-icon><UserFilled /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/system/roles">
            <el-icon><Key /></el-icon>
            <span>角色管理</span>
          </el-menu-item>
          <el-menu-item index="/system/menus">
            <el-icon><Menu /></el-icon>
            <span>菜单管理</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header>
        <div class="header-content">
          <span>欢迎，{{ userStore.user?.nickname || '管理员' }}</span>
          <el-button @click="handleLogout" type="danger" plain>退出</el-button>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import {
  DataBoard,
  TrendCharts,
  User,
  Notebook,
  Suitcase,
  Calendar,
  Warning,
  Wallet,
  Download,
  Avatar,
  House,
  Setting,
  UserFilled,
  Key,
  Menu as MenuIcon
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.el-aside {
  background-color: #304156;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #fff;
  background-color: #2a3c4f;
}

.el-header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
