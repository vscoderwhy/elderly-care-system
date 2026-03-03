<template>
  <view class="notifications-page">
    <!-- 顶部操作栏 -->
    <view class="page-header">
      <text class="header-title">消息通知</text>
      <view class="header-actions">
        <text class="read-all-btn" @click="handleReadAll">全部已读</text>
      </view>
    </view>

    <!-- 通知分类 -->
    <view class="category-tabs">
      <view
        class="tab-item"
        :class="{ active: activeTab === 'all' }"
        @click="activeTab = 'all'"
      >
        全部
      </view>
      <view
        class="tab-item"
        :class="{ active: activeTab === 'care' }"
        @click="activeTab = 'care'"
      >
        护理
      </view>
      <view
        class="tab-item"
        :class="{ active: activeTab === 'health' }"
        @click="activeTab = 'health'"
      >
        健康
      </view>
      <view
        class="tab-item"
        :class="{ active: activeTab === 'bill' }"
        @click="activeTab = 'bill'"
      >
        费用
      </view>
      <view
        class="tab-item"
        :class="{ active: activeTab === 'system' }"
        @click="activeTab = 'system'"
      >
        系统
      </view>
    </view>

    <!-- 通知列表 -->
    <scroll-view class="notifications-list" scroll-y @scrolltolower="loadMore">
      <view
        class="notification-item"
        :class="{ unread: !item.read }"
        v-for="item in notificationList"
        :key="item.id"
        @click="handleNotificationClick(item)"
      >
        <view class="item-left">
          <view class="notification-icon" :class="`type-${item.type}`">
            <uni-icons :name="getTypeIcon(item.type)" size="20" color="#fff" />
          </view>
          <view class="notification-content">
            <view class="content-header">
              <text class="content-title">{{ item.title }}</text>
              <text class="content-time">{{ formatTime(item.time) }}</text>
            </view>
            <text class="content-desc">{{ item.content }}</text>
          </view>
        </view>
        <view class="item-right">
          <uni-icons v-if="!item.read" type="checkbox" size="16" color="#409eff" />
        </view>
      </view>

      <!-- 加载更多 -->
      <view class="load-more" v-if="hasMore" @click="loadMore">
        <text>加载更多</text>
      </view>

      <!-- 空状态 -->
      <view v-if="!loading && notificationList.length === 0" class="empty-state">
        <text class="empty-text">暂无通知</text>
      </view>

      <!-- 加载中 -->
      <view v-if="loading" class="loading-state">
        <text>加载中...</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { formatRelativeTime } from '@/utils'

interface Notification {
  id: number
  type: string
  title: string
  content: string
  time: number
  read: boolean
  link?: string
}

const activeTab = ref('all')
const loading = ref(false)
const hasMore = ref(false)

const notificationList = ref<Notification[]>([
  {
    id: 1,
    type: 'care',
    title: '护理完成提醒',
    content: '张奶奶的日常护理已完成，服务人员：赵护士',
    time: Date.now() - 1800000,
    read: false,
    link: '/pages/care/detail?id=1'
  },
  {
    id: 2,
    type: 'health',
    title: '健康数据异常',
    content: '李爷爷血压偏高（150/95mmHg），请关注',
    time: Date.now() - 3600000,
    read: false,
    link: '/pages/elderly/health?id=2'
  },
  {
    id: 3,
    type: 'bill',
    title: '账单待支付提醒',
    content: '您有一笔账单待支付：床位费 ¥3,500.00',
    time: Date.now() - 86400000,
    read: true,
    link: '/pages/bills/detail?id=1'
  },
  {
    id: 4,
    type: 'care',
    title: '新护理记录',
    content: '王奶奶的康复训练已完成，训练时长30分钟',
    time: Date.now() - 172800000,
    read: true,
    link: '/pages/care/detail?id=2'
  },
  {
    id: 5,
    type: 'system',
    title: '系统通知',
    content: '系统将于3月5日02:00-04:00进行维护，届时将暂停服务',
    time: Date.now() - 259200000,
    read: true
  }
])

const formatTime = (timestamp: number) => {
  return formatRelativeTime(timestamp)
}

const getTypeIcon = (type: string) => {
  const map: Record<string, string> = {
    care: 'heart-filled',
    health: 'pulse',
    bill: 'wallet-filled',
    system: 'notification'
  }
  return map[type] || 'notification'
}

const handleNotificationClick = (item: Notification) => {
  // 标记为已读
  item.read = true

  // 跳转到详情页面
  if (item.link) {
    uni.navigateTo({
      url: item.link
    })
  }
}

const handleReadAll = () => {
  notificationList.value.forEach(item => {
    item.read = true
  })
  uni.showToast({
    title: '已全部标记为已读',
    icon: 'success'
  })
}

const loadMore = () => {
  if (!hasMore.value) return
  loading.value = true
  // TODO: 加载更多通知
  setTimeout(() => {
    loading.value = false
  }, 500)
}

onMounted(() => {
  // TODO: 加载通知列表
})
</script>

<style lang="scss" scoped>
.notifications-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  background: #fff;
}

.header-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #303133;
}

.read-all-btn {
  font-size: 26rpx;
  color: #409eff;
}

.category-tabs {
  display: flex;
  background: #fff;
  padding: 20rpx 30rpx;
  border-bottom: 1px solid #ebeef5;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  font-size: 28rpx;
  color: #606266;
  position: relative;
  transition: all 0.3s;

  &.active {
    color: #409eff;
    font-weight: 500;

    &::after {
      content: '';
      position: absolute;
      bottom: -20rpx;
      left: 50%;
      transform: translateX(-50%);
      width: 40rpx;
      height: 4rpx;
      background: #409eff;
      border-radius: 2rpx;
    }
  }
}

.notifications-list {
  padding: 20rpx 30rpx;
  height: calc(100vh - 180rpx);
}

.notification-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16rpx;
  padding: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  transition: all 0.3s;

  &.unread {
    background: #ecf5ff;

    .notification-icon {
      position: relative;

      &::after {
        content: '';
        position: absolute;
        top: -4rpx;
        right: -4rpx;
        width: 12rpx;
        height: 12rpx;
        background: #f56c6c;
        border-radius: 50%;
        border: 2rpx solid #fff;
      }
    }
  }
}

.item-left {
  flex: 1;
  display: flex;
  gap: 16rpx;
}

.notification-icon {
  width: 64rpx;
  height: 64rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &.type-care { background: #67c23a; }
  &.type-health { background: #e6a23c; }
  &.type-bill { background: #f56c6c; }
  &.type-system { background: #909399; }
}

.notification-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.content-title {
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;
}

.content-time {
  font-size: 24rpx;
  color: #909399;
}

.content-desc {
  font-size: 26rpx;
  color: #606266;
  line-height: 1.5;
}

.item-right {
  flex-shrink: 0;
}

.load-more {
  padding: 30rpx;
  text-align: center;
  font-size: 28rpx;
  color: #409eff;
}

.empty-state {
  padding: 120rpx 0;
  text-align: center;
}

.empty-text {
  font-size: 28rpx;
  color: #909399;
}

.loading-state {
  padding: 80rpx 0;
  text-align: center;
  font-size: 28rpx;
  color: #909399;
}
</style>
