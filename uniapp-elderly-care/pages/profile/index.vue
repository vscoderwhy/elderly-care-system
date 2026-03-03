<template>
  <view class="profile-page">
    <!-- 用户信息卡片 -->
    <view class="user-card">
      <view class="user-info" @click="handleEditProfile">
        <image class="user-avatar" :src="userInfo.avatar || '/static/avatar-default.png'" mode="aspectFill" />
        <view class="user-detail">
          <text class="user-name">{{ userInfo.name }}</text>
          <text class="user-role">{{ getRoleText(userInfo.role) }}</text>
        </view>
        <view class="edit-icon">
          <uni-icons type="compose" size="18" color="#909399" />
        </view>
      </view>
    </view>

    <!-- 功能列表 -->
    <view class="menu-section">
      <view class="menu-group">
        <view class="menu-item" @click="handleNavigate('elderly')">
          <view class="item-left">
            <view class="item-icon" style="background: #e1f3e8">
              <uni-icons type="person" size="20" color="#67c23a" />
            </view>
            <text class="item-label">我的老人</text>
          </view>
          <view class="item-right">
            <text class="item-value">{{ elderlyCount }}位</text>
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>

        <view class="menu-item" @click="handleNavigate('visits')">
          <view class="item-left">
            <view class="item-icon" style="background: #ecf5ff">
              <uni-icons type="calendar" size="20" color="#409eff" />
            </view>
            <text class="item-label">预约记录</text>
          </view>
          <view class="item-right">
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>

        <view class="menu-item" @click="handleNavigate('bills')">
          <view class="item-left">
            <view class="item-icon" style="background: #fef0f0">
              <uni-icons type="wallet-filled" size="20" color="#f56c6c" />
            </view>
            <text class="item-label">账单记录</text>
          </view>
          <view class="item-right">
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>
      </view>

      <view class="menu-group">
        <view class="menu-item" @click="handleNavigate('settings')">
          <view class="item-left">
            <view class="item-icon" style="background: #fdf6ec">
              <uni-icons type="gear" size="20" color="#e6a23c" />
            </view>
            <text class="item-label">设置</text>
          </view>
          <view class="item-right">
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>

        <view class="menu-item" @click="handleNavigate('about')">
          <view class="item-left">
            <view class="item-icon" style="background: #f4f4f5">
              <uni-icons type="info" size="20" color="#909399" />
            </view>
            <text class="item-label">关于</text>
          </view>
          <view class="item-right">
            <text class="item-value">v1.0.0</text>
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>
      </view>

      <view class="menu-group">
        <view class="menu-item" @click="handleContact">
          <view class="item-left">
            <view class="item-icon" style="background: #e1f3ff">
              <uni-icons type="phone" size="20" color="#409eff" />
            </view>
            <text class="item-label">联系客服</text>
          </view>
          <view class="item-right">
            <text class="item-value">400-123-4567</text>
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>

        <view class="menu-item" @click="handleFeedback">
          <view class="item-left">
            <view class="item-icon" style="background: #fef9f3">
              <uni-icons type="chat" size="20" color="#e6a23c" />
            </view>
            <text class="item-label">意见反馈</text>
          </view>
          <view class="item-right">
            <uni-icons type="arrowright" size="14" color="#909399" />
          </view>
        </view>
      </view>
    </view>

    <!-- 退出登录 -->
    <view class="logout-section">
      <button class="logout-btn" @click="handleLogout">退出登录</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

const userInfo = ref({
  id: '',
  name: '家属用户',
  avatar: '',
  role: 'family',
  phone: ''
})

const elderlyCount = ref(2)

const getRoleText = (role: string) => {
  const map: Record<string, string> = {
    family: '家属',
    caregiver: '护工',
    admin: '管理员'
  }
  return map[role] || role
}

const handleEditProfile = () => {
  uni.navigateTo({
    url: '/pages/profile/edit'
  })
}

const handleNavigate = (page: string) => {
  const routes: Record<string, string> = {
    elderly: '/pages/profile/elderly',
    visits: '/pages/profile/visits',
    bills: '/pages/profile/bills',
    settings: '/pages/profile/settings',
    about: '/pages/profile/about'
  }

  if (routes[page]) {
    if (page === 'settings') {
      uni.navigateTo({ url: routes[page] })
    } else if (page === 'about') {
      uni.navigateTo({ url: routes[page] })
    } else {
      uni.navigateTo({ url: routes[page] })
    }
  }
}

const handleContact = () => {
  uni.showModal({
    title: '联系客服',
    content: '客服电话：400-123-4567\n工作时间：8:00-18:00',
    confirmText: '拨打电话',
    success: (res) => {
      if (res.confirm) {
        uni.makePhoneCall({
          phoneNumber: '4001234567'
        })
      }
    }
  })
}

const handleFeedback = () => {
  uni.navigateTo({
    url: '/pages/profile/feedback'
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
  // 加载用户信息
  const user = userStore.user
  if (user) {
    userInfo.value = {
      id: user.id,
      name: user.name,
      avatar: user.avatar,
      role: user.role,
      phone: user.phone
    }
  }
})
</script>

<style lang="scss" scoped>
.profile-page {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 40rpx;
}

.user-card {
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  padding: 60rpx 30rpx 40rpx;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.user-avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
}

.user-detail {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.user-name {
  font-size: 36rpx;
  font-weight: 600;
  color: #fff;
}

.user-role {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
}

.edit-icon {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
}

.menu-section {
  padding: 20rpx 30rpx;
}

.menu-group {
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 20rpx;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx 24rpx;
  border-bottom: 1px solid #f5f7fa;

  &:last-child {
    border-bottom: none;
  }
}

.item-left {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.item-icon {
  width: 56rpx;
  height: 56rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.item-label {
  font-size: 28rpx;
  color: #303133;
}

.item-right {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.item-value {
  font-size: 26rpx;
  color: #909399;
}

.logout-section {
  padding: 0 30rpx;
  margin-top: 40rpx;
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
