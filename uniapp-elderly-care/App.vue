<template>
  <view class="app">
    <view class="app-content">
      <!-- 页面内容 -->
    </view>
  </view>
</template>

<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

onLaunch(() => {
  console.log('App Launch')
  // 初始化用户状态
  userStore.init()

  // 检查更新
  checkUpdate()
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})

// 检查更新（小程序）
const checkUpdate = () => {
  // #ifdef MP-WEIXIN
  const updateManager = uni.getUpdateManager()

  updateManager.onCheckForUpdate((res) => {
    console.log('检查更新结果:', res.hasUpdate)
  })

  updateManager.onUpdateReady(() => {
    uni.showModal({
      title: '更新提示',
      content: '新版本已经准备好，是否重启应用？',
      success: (res) => {
        if (res.confirm) {
          updateManager.applyUpdate()
        }
      }
    })
  })

  updateManager.onUpdateFailed(() => {
    uni.showModal({
      title: '更新失败',
      content: '新版本下载失败，请检查网络后重试',
      showCancel: false
    })
  })
  // #endif
}
</script>

<style lang="scss">
/* 全局样式 */
page {
  background-color: #f5f7fa;
  font-size: 14px;
  line-height: 1.5;
}

.app {
  width: 100%;
  min-height: 100vh;
}

.app-content {
  width: 100%;
  min-height: 100vh;
}

/* 全局工具类 */
.flex {
  display: flex;
}

.flex-row {
  display: flex;
  flex-direction: row;
}

.flex-column {
  display: flex;
  flex-direction: column;
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

.flex-around {
  display: flex;
  align-items: center;
  justify-content: space-around;
}

.flex-1 {
  flex: 1;
}

/* 文本对齐 */
.text-left {
  text-align: left;
}

.text-center {
  text-align: center;
}

.text-right {
  text-align: right;
}

/* 文本颜色 */
.text-primary {
  color: #409eff;
}

.text-success {
  color: #67c23a;
}

.text-warning {
  color: #e6a23c;
}

.text-danger {
  color: #f56c6c;
}

.text-info {
  color: #909399;
}

.text-muted {
  color: #909399;
}

/* 字体大小 */
.text-xs {
  font-size: 10px;
}

.text-sm {
  font-size: 12px;
}

.text-base {
  font-size: 14px;
}

.text-lg {
  font-size: 16px;
}

.text-xl {
  font-size: 18px;
}

.text-2xl {
  font-size: 20px;
}

/* 字体粗细 */
.font-normal {
  font-weight: 400;
}

.font-medium {
  font-weight: 500;
}

.font-bold {
  font-weight: 700;
}

/* 间距 */
.m-0 {
  margin: 0;
}

.m-10 {
  margin: 10px;
}

.m-20 {
  margin: 20px;
}

.mt-5 {
  margin-top: 5px;
}

.mt-10 {
  margin-top: 10px;
}

.mt-20 {
  margin-top: 20px;
}

.mb-5 {
  margin-bottom: 5px;
}

.mb-10 {
  margin-bottom: 10px;
}

.mb-20 {
  margin-bottom: 20px;
}

.ml-5 {
  margin-left: 5px;
}

.ml-10 {
  margin-left: 10px;
}

.mr-5 {
  margin-right: 5px;
}

.mr-10 {
  margin-right: 10px;
}

.p-0 {
  padding: 0;
}

.p-10 {
  padding: 10px;
}

.p-20 {
  padding: 20px;
}

.pt-5 {
  padding-top: 5px;
}

.pt-10 {
  padding-top: 10px;
}

.pb-5 {
  padding-bottom: 5px;
}

.pb-10 {
  padding-bottom: 10px;
}

.pl-10 {
  padding-left: 10px;
}

.pr-10 {
  padding-right: 10px;
}

/* 宽高 */
.w-full {
  width: 100%;
}

.h-full {
  height: 100%;
}

/* 边框 */
.border {
  border: 1px solid #dcdfe6;
}

.border-b {
  border-bottom: 1px solid #dcdfe6;
}

.border-t {
  border-top: 1px solid #dcdfe6;
}

.border-l {
  border-left: 1px solid #dcdfe6;
}

.border-r {
  border-right: 1px solid #dcdfe6;
}

/* 圆角 */
.rounded {
  border-radius: 4px;
}

.rounded-full {
  border-radius: 9999px;
}

/* 阴影 */
.shadow {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.shadow-lg {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 溢出 */
.overflow-hidden {
  overflow: hidden;
}

.overflow-auto {
  overflow: auto;
}

/* 定位 */
.relative {
  position: relative;
}

.absolute {
  position: absolute;
}

.fixed {
  position: fixed;
}

/* 显示隐藏 */
.hidden {
  display: none;
}

.visible {
  visibility: visible;
}

.invisible {
  visibility: hidden;
}

/* 卡片样式 */
.card {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.card-body {
  padding: 0;
}

/* 列表项样式 */
.list-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;

  &:last-child {
    border-bottom: none;
  }
}

.list-item-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 12px;
  flex-shrink: 0;
}

.list-item-content {
  flex: 1;
  overflow: hidden;
}

.list-item-title {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
}

.list-item-desc {
  font-size: 12px;
  color: #909399;
}

.list-item-extra {
  flex-shrink: 0;
  margin-left: 12px;
}

/* 标签样式 */
.tag {
  display: inline-block;
  padding: 2px 8px;
  font-size: 12px;
  border-radius: 4px;
  background: #f0f2f5;
  color: #606266;
}

.tag-primary {
  background: #ecf5ff;
  color: #409eff;
}

.tag-success {
  background: #f0f9ff;
  color: #67c23a;
}

.tag-warning {
  background: #fdf6ec;
  color: #e6a23c;
}

.tag-danger {
  background: #fef0f0;
  color: #f56c6c;
}

/* 徽标样式 */
.badge {
  position: relative;
  display: inline-block;
}

.badge-count {
  position: absolute;
  top: -8px;
  right: -8px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  border-radius: 8px;
  background: #f56c6c;
  color: #fff;
  font-size: 10px;
  line-height: 16px;
  text-align: center;
}

/* 安全区域适配 */
.safe-area-inset-top {
  padding-top: constant(safe-area-inset-top);
  padding-top: env(safe-area-inset-top);
}

.safe-area-inset-bottom {
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}
</style>
