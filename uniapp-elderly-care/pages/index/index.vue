<template>
  <view class="index-page">
    <!-- 头部区域 -->
    <view class="page-header">
      <view class="header-top">
        <view class="user-info" @click="handleSelectElderly">
          <view class="avatar-wrapper">
            <image v-if="currentElderly.avatar" :src="currentElderly.avatar" class="avatar" />
            <view v-else class="avatar-placeholder">
              <text class="avatar-text">{{ currentElderly.name?.charAt(0) || '选' }}</text>
            </view>
          </view>
          <view class="user-detail">
            <text class="user-name">{{ currentElderly.name || '请选择老人' }}</text>
            <text class="user-relation" v-if="currentElderly.relation">
              {{ currentElderly.relation }}
            </text>
          </view>
          <text class="change-text">{{ elderlyList.length > 1 ? '切换 >' : '' }}</text>
        </view>
        <view class="header-actions">
          <view class="action-item" @click="handleScan">
            <uni-icons type="scan" size="24" color="#333" />
          </view>
          <view class="action-item" @click="handleNotice">
            <uni-icons type="notification" size="24" color="#333" />
            <view v-if="unreadCount > 0" class="badge">{{ unreadCount }}</view>
          </view>
        </view>
      </view>
    </view>

    <!-- 统计卡片 -->
    <view class="stats-section">
      <view class="stat-card" v-for="stat in stats" :key="stat.key" @click="handleStatClick(stat)">
        <view class="stat-icon" :class="`stat-${stat.type}`">
          <uni-icons :type="stat.icon" size="28" color="#fff" />
        </view>
        <view class="stat-content">
          <text class="stat-value">{{ stat.value }}</text>
          <text class="stat-label">{{ stat.label }}</text>
        </view>
      </view>
    </view>

    <!-- 快捷功能 -->
    <view class="quick-actions">
      <view class="action-title">快捷功能</view>
      <view class="action-grid">
        <view class="action-item" v-for="action in quickActions" :key="action.key" @click="handleActionClick(action)">
          <view class="action-icon" :style="{ background: action.color }">
            <uni-icons :type="action.icon" size="24" color="#fff" />
          </view>
          <text class="action-text">{{ action.label }}</text>
        </view>
      </view>
    </view>

    <!-- 健康数据卡片 -->
    <view class="health-card" @click="handleHealthDetail">
      <view class="card-header">
        <text class="card-title">今日健康</text>
        <text class="card-more">详情 ></text>
      </view>
      <view class="health-data">
        <view class="health-item" v-for="item in healthData" :key="item.key">
          <text class="health-label">{{ item.label }}</text>
          <view class="health-value-row">
            <text class="health-value" :class="`health-${item.level}`">{{ item.value }}</text>
            <text class="health-unit">{{ item.unit }}</text>
            <view class="health-trend" :class="item.trend >= 0 ? 'up' : 'down'">
              <uni-icons :type="item.trend >= 0 ? 'arrow-up' : 'arrow-down'" size="12" />
              <text>{{ Math.abs(item.trend) }}%</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- 最近护理记录 -->
    <view class="care-records">
      <view class="section-header">
        <text class="section-title">最近护理</text>
        <text class="section-more" @click="handleMoreCare">查看全部 ></text>
      </view>
      <view class="record-list" v-if="recentRecords.length > 0">
        <view class="record-item" v-for="record in recentRecords" :key="record.id" @click="handleRecordClick(record)">
          <view class="record-icon">
            <uni-icons type="heart" size="20" color="#409eff" />
          </view>
          <view class="record-content">
            <view class="record-header-row">
              <text class="record-title">{{ record.type }}</text>
              <text class="record-time">{{ formatTime(record.time) }}</text>
            </view>
            <text class="record-desc">{{ record.content }}</text>
          </view>
        </view>
      </view>
      <view v-else class="empty-state">
        <text>暂无护理记录</text>
      </view>
    </view>

    <!-- 老人选择弹窗 -->
    <uni-popup ref="elderlyPopup" type="bottom">
      <view class="elderly-selector">
        <view class="selector-header">
          <text class="selector-title">选择老人</text>
          <uni-icons type="close" size="20" @click="closeElderlyPopup" />
        </view>
        <scroll-view class="selector-list" scroll-y>
          <view
            class="selector-item"
            v-for="elderly in elderlyList"
            :key="elderly.id"
            @click="selectElderly(elderly)"
          >
            <image v-if="elderly.avatar" :src="elderly.avatar" class="selector-avatar" />
            <view v-else class="selector-avatar-placeholder">
              <text>{{ elderly.name?.charAt(0) }}</text>
            </view>
            <view class="selector-info">
              <text class="selector-name">{{ elderly.name }}</text>
              <text class="selector-relation">{{ elderly.relation }}</text>
            </view>
            <uni-icons v-if="currentElderly.id === elderly.id" type="checkmarkempty" size="20" color="#409eff" />
          </view>
        </scroll-view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { formatRelativeTime } from '@/utils'

const userStore = useUserStore()

// 当前选中的老人
const currentElderly = ref<any>({
  id: '',
  name: '',
  avatar: '',
  relation: ''
})

// 老人列表
const elderlyList = ref<any[]>([])

// 未读消息数
const unreadCount = ref(5)

// 统计数据
const stats = ref([
  { key: 'care', label: '本月护理', value: '128次', type: 'primary', icon: 'heart-filled' },
  { key: 'health', label: '健康记录', value: '56次', type: 'success', icon: 'checkmarkempty' },
  { key: 'bill', label: '待付账单', value: '3笔', type: 'warning', icon: 'wallet' },
  { key: 'visit', label: '探视记录', value: '8次', type: 'info', icon: 'calendar' }
])

// 快捷功能
const quickActions = ref([
  { key: 'care', label: '护理记录', icon: 'list', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { key: 'health', label: '健康数据', icon: 'pulse', color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { key: 'bill', label: '费用账单', icon: 'wallet-filled', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
  { key: 'visit', label: '探视预约', icon: 'calendar-filled', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' },
  { key: 'message', label: '消息通知', icon: 'chatbubble-filled', color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' },
  { key: 'contact', label: '联系护工', icon: 'phone-filled', color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)' },
  { key: 'photo', label: '照片相册', icon: 'images', color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)' },
  { key: 'more', label: '更多', icon: 'bars', color: 'linear-gradient(135deg, #d299c2 0%, #fef9d7 100%)' }
])

// 健康数据
const healthData = ref([
  { key: 'bp', label: '血压', value: '128/82', unit: 'mmHg', level: 'normal', trend: -2 },
  { key: 'hr', label: '心率', value: '72', unit: '次/分', level: 'normal', trend: 0 },
  { key: 'bs', label: '血糖', value: '6.3', unit: 'mmol/L', level: 'normal', trend: 5 },
  { key: 'temp', label: '体温', value: '36.5', unit: '℃', level: 'normal', trend: 0 }
])

// 最近护理记录
const recentRecords = ref([
  {
    id: 1,
    type: '日常护理',
    content: '测量血压、体温，协助用餐',
    time: Date.now() - 3600000,
    nurseName: '赵护士'
  },
  {
    id: 2,
    type: '健康监测',
    content: '血压130/85mmHg，心率72次/分',
    time: Date.now() - 7200000,
    nurseName: '周护士'
  },
  {
    id: 3,
    type: '康复训练',
    content: '上肢关节活动训练30分钟',
    time: Date.now() - 86400000,
    nurseName: '陈康复师'
  }
])

// 选择老人
const handleSelectElderly = () => {
  if (elderlyList.value.length > 1) {
    // 打开选择弹窗
    uni.$emit('openElderlySelector')
  }
}

const closeElderlyPopup = () => {
  uni.$emit('closeElderlySelector')
}

const selectElderly = (elderly: any) => {
  currentElderly.value = elderly
  userStore.setSelectedElderly(elderly.id)
  closeElderlyPopup()
  // 刷新数据
  loadData()
}

// 扫码
const handleScan = () => {
  uni.scanCode({
    success: (res) => {
      console.log('扫码结果', res)
    }
  })
}

// 消息通知
const handleNotice = () => {
  uni.navigateTo({
    url: '/pages/profile/notifications'
  })
}

// 统计卡片点击
const handleStatClick = (stat: any) => {
  console.log('统计卡片点击', stat)
}

// 快捷功能点击
const handleActionClick = (action: any) => {
  const routes: Record<string, string> = {
    care: '/pages/care/records',
    health: '/pages/elderly/health',
    bill: '/pages/bills/list',
    visit: '/pages/visits/appointment',
    message: '/pages/profile/notifications',
    contact: '/pages/profile/contact',
    photo: '/pages/elderly/photos'
  }

  if (routes[action.key]) {
    uni.navigateTo({ url: routes[action.key] })
  } else if (action.key === 'more') {
    uni.showActionSheet({
      itemList: ['设置', '关于', '退出登录'],
      success: (res) => {
        if (res.tapIndex === 0) {
          uni.navigateTo({ url: '/pages/profile/settings' })
        } else if (res.tapIndex === 2) {
          userStore.logout()
        }
      }
    })
  }
}

// 健康详情
const handleHealthDetail = () => {
  uni.navigateTo({
    url: `/pages/elderly/health?id=${currentElderly.value.id}`
  })
}

// 查看更多护理记录
const handleMoreCare = () => {
  uni.navigateTo({
    url: '/pages/care/records'
  })
}

// 护理记录点击
const handleRecordClick = (record: any) => {
  console.log('记录详情', record)
}

// 格式化时间
const formatTime = (timestamp: number) => {
  return formatRelativeTime(timestamp)
}

// 加载数据
const loadData = async () => {
  try {
    // 获取关联老人列表
    // elderlyList.value = await getElderlyList()

    // 模拟数据
    elderlyList.value = [
      { id: 1, name: '张奶奶', relation: '母子', avatar: '' },
      { id: 2, name: '李爷爷', relation: '父女', avatar: '' }
    ]

    // 设置当前老人
    const selectedId = userStore.getSelectedElderly()
    if (selectedId) {
      currentElderly.value = elderlyList.value.find(e => e.id === selectedId) || elderlyList.value[0]
    } else if (elderlyList.value.length > 0) {
      currentElderly.value = elderlyList.value[0]
      userStore.setSelectedElderly(elderlyList.value[0].id)
    }
  } catch (error) {
    console.error('加载数据失败', error)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
.index-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20rpx 30rpx 30rpx;
  border-radius: 0 0 40rpx 40rpx;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 20rpx;
  flex: 1;
}

.avatar-wrapper {
  position: relative;
}

.avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
}

.avatar-placeholder {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
}

.avatar-text {
  font-size: 36rpx;
  color: #fff;
  font-weight: 600;
}

.user-detail {
  flex: 1;
}

.user-name {
  display: block;
  font-size: 32rpx;
  font-weight: 600;
  color: #fff;
  margin-bottom: 8rpx;
}

.user-relation {
  display: block;
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
}

.change-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
}

.header-actions {
  display: flex;
  gap: 20rpx;
}

.action-item {
  position: relative;
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.badge {
  position: absolute;
  top: -4rpx;
  right: -4rpx;
  min-width: 32rpx;
  height: 32rpx;
  padding: 0 8rpx;
  background: #f56c6c;
  border-radius: 16rpx;
  font-size: 20rpx;
  color: #fff;
  line-height: 32rpx;
  text-align: center;
}

.stats-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
  padding: 30rpx;
  margin-top: -60rpx;
}

.stat-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  display: flex;
  align-items: center;
  gap: 20rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
}

.stat-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;

  &.stat-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.stat-success {
    background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
  }

  &.stat-warning {
    background: linear-gradient(135deg, #fccb90 0%, #d57eeb 100%);
  }

  &.stat-info {
    background: linear-gradient(135deg, #e0c3fc 0%, #8ec5fc 100%);
  }
}

.stat-content {
  flex: 1;
}

.stat-value {
  display: block;
  font-size: 32rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4rpx;
}

.stat-label {
  display: block;
  font-size: 24rpx;
  color: #909399;
}

.quick-actions {
  padding: 0 30rpx 30rpx;
}

.action-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20rpx;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20rpx;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
}

.action-icon {
  width: 100rpx;
  height: 100rpx;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
}

.action-text {
  font-size: 24rpx;
  color: #606266;
}

.health-card {
  margin: 0 30rpx 30rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.card-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
}

.card-more {
  font-size: 24rpx;
  color: #909399;
}

.health-data {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.health-item {
  padding: 20rpx;
  background: #f5f7fa;
  border-radius: 12rpx;
}

.health-label {
  display: block;
  font-size: 24rpx;
  color: #909399;
  margin-bottom: 8rpx;
}

.health-value-row {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
}

.health-value {
  font-size: 36rpx;
  font-weight: 600;
  color: #303133;

  &.health-normal {
    color: #67c23a;
  }

  &.health-warning {
    color: #e6a23c;
  }

  &.health-danger {
    color: #f56c6c;
  }
}

.health-unit {
  font-size: 24rpx;
  color: #909399;
}

.health-trend {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 4rpx;
  font-size: 20rpx;

  &.up {
    color: #f56c6c;
  }

  &.down {
    color: #67c23a;
  }
}

.care-records {
  padding: 0 30rpx 30rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
}

.section-more {
  font-size: 24rpx;
  color: #909399;
}

.record-list {
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}

.record-item {
  display: flex;
  gap: 20rpx;
  padding: 24rpx;
  border-bottom: 1rpx solid #ebeef5;

  &:last-child {
    border-bottom: none;
  }
}

.record-icon {
  width: 60rpx;
  height: 60rpx;
  border-radius: 12rpx;
  background: #ecf5ff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.record-content {
  flex: 1;
}

.record-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8rpx;
}

.record-title {
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;
}

.record-time {
  font-size: 24rpx;
  color: #909399;
}

.record-desc {
  font-size: 26rpx;
  color: #606266;
  line-height: 1.5;
}

.empty-state {
  padding: 80rpx 0;
  text-align: center;
  font-size: 28rpx;
  color: #909399;
}

// 老人选择弹窗
.elderly-selector {
  background: #fff;
  border-radius: 32rpx 32rpx 0 0;
  padding: 30rpx 0 0;
  max-height: 80vh;
}

.selector-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 30rpx 30rpx;
  border-bottom: 1rpx solid #ebeef5;
}

.selector-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #303133;
}

.selector-list {
  max-height: 60vh;
}

.selector-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx 30rpx;
  border-bottom: 1rpx solid #f5f7fa;

  &:last-child {
    border-bottom: none;
  }
}

.selector-avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
}

.selector-avatar-placeholder {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  color: #909399;
}

.selector-info {
  flex: 1;
}

.selector-name {
  display: block;
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4rpx;
}

.selector-relation {
  display: block;
  font-size: 24rpx;
  color: #909399;
}
</style>
