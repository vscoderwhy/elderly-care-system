<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useDarkMode } from '@/composables/useDarkMode'

// 初始化暗黑模式
const { updateTheme } = useDarkMode()

onMounted(() => {
  updateTheme()
})
</script>

<style lang="scss">
// 全局样式
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  width: 100%;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-primary);
  background-color: var(--bg-secondary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  width: 100%;
  height: 100%;
}

// 滚动条样式
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-thumb {
  background: var(--border-color-light);
  border-radius: 3px;

  &:hover {
    background: var(--border-color);
  }
}

::-webkit-scrollbar-track {
  background: transparent;
}

// 链接样式
a {
  color: var(--primary-color);
  text-decoration: none;
  transition: color 0.3s;

  &:hover {
    color: var(--primary-color);
    opacity: 0.8;
  }
}

// 过渡动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// 工具类
.text-center {
  text-align: center;
}

.text-left {
  text-align: left;
}

.text-right {
  text-align: right;
}

.mt-10 {
  margin-top: 10px;
}

.mt-20 {
  margin-top: 20px;
}

.mb-10 {
  margin-bottom: 10px;
}

.mb-20 {
  margin-bottom: 20px;
}

.p-10 {
  padding: 10px;
}

.p-20 {
  padding: 20px;
}

.flex {
  display: flex;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.full-width {
  width: 100%;
}

.full-height {
  height: 100%;
}

// Element Plus 样式覆盖
.el-card {
  --el-card-border-color: var(--border-color);
}

.el-table {
  --el-table-border-color: var(--border-color-lighter);
  --el-table-bg-color: var(--card-bg);
  --el-table-tr-bg-color: var(--card-bg);
  --el-table-current-row-bg-color: var(--bg-tertiary);
  --el-table-header-bg-color: var(--bg-secondary);
}

.el-menu {
  --el-menu-bg-color: var(--bg-secondary);
  --el-menu-text-color: var(--text-primary);
  --el-menu-hover-bg-color: var(--bg-tertiary);
  --el-menu-active-color: var(--primary-color);
}

.el-input__wrapper,
.el-textarea__inner {
  background-color: var(--bg-tertiary);
}

// 响应式断点
@media (max-width: 768px) {
  html {
    font-size: 12px;
  }

  .el-message-box {
    width: 90% !important;
  }

  .el-dialog {
    width: 95% !important;
  }
}
</style>
