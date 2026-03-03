<template>
  <view class="settings-page">
    <!-- 顶部状态栏 -->
    <view class="status-bar">
      <view class="status-back" @click="handleBack">
        <uni-icons type="back" size="20" color="#606266" />
      </view>
      <text class="status-title">设置</text>
      <view class="status-placeholder"></view>
    </view>

    <scroll-view class="settings-content" scroll-y>
      <!-- 账号设置 -->
      <view class="setting-group">
        <view class="group-title">账号设置</view>
        <view class="setting-item" @click="handleNavigate('phone')">
          <view class="item-left">
            <text class="item-label">手机号</text>
          </view>
          <view class="item-right">
            <text class="item-value">{{ userInfo.phone || '未绑定' }}</text>
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>
        <view class="setting-item" @click="handleNavigate('password')">
          <view class="item-left">
            <text class="item-label">登录密码</text>
          </view>
          <view class="item-right">
            <text class="item-value">已设置</text>
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>
      </view>

      <!-- 通知设置 -->
      <view class="setting-group">
        <view class="group-title">通知设置</view>
        <view class="setting-item">
          <view class="item-left">
            <text class="item-label">护理通知</text>
          </view>
          <view class="item-right">
            <switch :checked="settings.careNotification" @change="handleSettingChange('careNotification', $event)" color="#409eff" />
          </view>
        </view>
        <view class="setting-item">
          <view class="item-left">
            <text class="item-label">健康提醒</text>
          </view>
          <view class="item-right">
            <switch :checked="settings.healthReminder" @change="handleSettingChange('healthReminder', $event)" color="#409eff" />
          </view>
        </view>
        <view class="setting-item">
          <view class="item-left">
            <text class="item-label">账单提醒</text>
          </view>
          <view class="item-right">
            <switch :checked="settings.billReminder" @change="handleSettingChange('billReminder', $event)" color="#409eff" />
          </view>
        </view>
        <view class="setting-item">
          <view class="item-left">
            <text class="item-label">系统通知</text>
          </view>
          <view class="item-right">
            <switch :checked="settings.systemNotification" @change="handleSettingChange('systemNotification', $event)" color="#409eff" />
          </view>
        </view>
      </view>

      <!-- 通用设置 -->
      <view class="setting-group">
        <view class="group-title">通用</view>
        <view class="setting-item" @click="handleNavigate('language')">
          <view class="item-left">
            <text class="item-label">语言</text>
          </view>
          <view class="item-right">
            <text class="item-value">简体中文</text>
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>
        <view class="setting-item" @click="handleClearCache">
          <view class="item-left">
            <text class="item-label">清除缓存</text>
          </view>
          <view class="item-right">
            <text class="item-value">{{ cacheSize }}</text>
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>
        <view class="setting-item" @click="handleNavigate('about')">
          <view class="item-left">
            <text class="item-label">关于我们</text>
          </view>
          <view class="item-right">
            <text class="item-value">v1.0.0</text>
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>
      </view>

      <!-- 退出登录 -->
      <view class="logout-section">
        <button class="logout-btn" @click="handleLogout">退出登录</button>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

const userInfo = ref({
  phone: '138****8888'
})

const settings = ref({
  careNotification: true,
  healthReminder: true,
  billReminder: true,
  systemNotification: true
})

const cacheSize = ref('23.5MB')

const handleBack = () => {
  uni.navigateBack()
}

const handleNavigate = (page: string) => {
  console.log('导航到:', page)
}

const handleSettingChange = (key: string, e: any) => {
  settings.value[key] = e.detail.value
  uni.showToast({
    title: '设置已更新',
    icon: 'success'
  })
  // TODO: 保存设置到服务器
}

const handleClearCache = () => {
  uni.showModal({
    title: '提示',
    content: '确定要清除缓存吗？',
    success: (res) => {
      if (res.confirm) {
        // TODO: 清除缓存逻辑
        cacheSize.value = '0KB'
        uni.showToast({
          title: '缓存已清除',
          icon: 'success'
        })
      }
    }
  })
}

const handleLogout = () => {
  uni.showModal({
    title: '提示',
    content: '确定要退出登录吗？',
    success: async (res) => {
      if (res.confirm) {
        await userStore.logout()
        uni.reLaunch({
          url: '/pages/login/index'
        })
      }
    }
  })
}

onMounted(() => {
  // TODO: 加载用户设置
})
</script>

<style lang="scss" scoped>
.settings-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  background: #fff;
  position: sticky;
  top: 0;
  z-index: 100;
}

.status-back,
.status-placeholder {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #303133;
}

.settings-content {
  padding: 20rpx 0;
}

.setting-group {
  background: #fff;
  margin-bottom: 20rpx;
}

.group-title {
  padding: 30rpx 30rpx 20rpx;
  font-size: 26rpx;
  color: #909399;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1px solid #f5f7fa;

  &:last-child {
    border-bottom: none;
  }
}

.item-left {
  flex: 1;
}

.item-label {
  font-size: 28rpx;
  color: #303133;
}

.item-right {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.item-value {
  font-size: 26rpx;
  color: #909399;
}

.logout-section {
  padding: 40rpx 30rpx;
}

.logout-btn {
  width: 100%;
  height: 80rpx;
  background: #fff;
  color: #f56c6c;
  border: none;
  border-radius: 40rpx;
  font-size: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
