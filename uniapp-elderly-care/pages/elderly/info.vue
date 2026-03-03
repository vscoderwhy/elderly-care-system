<template>
  <view class="elderly-info-page">
    <!-- 头部卡片 -->
    <view class="elderly-card">
      <view class="card-bg"></view>
      <view class="card-content">
        <view class="elderly-avatar-wrapper">
          <image v-if="elderlyInfo.avatar" :src="elderlyInfo.avatar" class="elderly-avatar" mode="aspectFill" />
          <view v-else class="elderly-avatar-placeholder">
            <text class="avatar-text">{{ elderlyInfo.name?.charAt(0) || '老' }}</text>
          </view>
          <view class="care-level-tag" :class="`level-${elderlyInfo.careLevel}`">
            {{ getCareLevelText(elderlyInfo.careLevel) }}
          </view>
        </view>

        <view class="elderly-base-info">
          <text class="elderly-name">{{ elderlyInfo.name }}</text>
          <view class="info-tags">
            <view class="info-tag">
              <text class="tag-text">{{ elderlyInfo.gender }}</text>
            </view>
            <view class="info-tag">
              <text class="tag-text">{{ elderlyInfo.age }}岁</text>
            </view>
            <view class="info-tag">
              <text class="tag-text">{{ elderlyInfo.bedNumber }}</text>
            </view>
          </view>
        </view>

        <view class="elderly-stats">
          <view class="stat-item">
            <text class="stat-value">{{ elderlyInfo.stayDays }}</text>
            <text class="stat-label">入住天数</text>
          </view>
          <view class="stat-divider"></view>
          <view class="stat-item">
            <text class="stat-value">{{ elderlyInfo.careCount }}</text>
            <text class="stat-label">本月护理</text>
          </view>
          <view class="stat-divider"></view>
          <view class="stat-item">
            <text class="stat-value">{{ elderlyInfo.healthScore }}</text>
            <text class="stat-label">健康评分</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 标签页 -->
    <view class="tabs-container">
      <view class="tabs-header">
        <view
          v-for="tab in tabs"
          :key="tab.key"
          class="tab-item"
          :class="{ active: activeTab === tab.key }"
          @click="activeTab = tab.key"
        >
          <text class="tab-text">{{ tab.label }}</text>
          <view v-if="activeTab === tab.key" class="tab-indicator"></view>
        </view>
      </view>

      <scroll-view class="tabs-content" scroll-y>
        <!-- 基本信息 -->
        <view v-if="activeTab === 'info'" class="tab-content">
          <view class="info-section">
            <view class="section-title">基本信息</view>
            <view class="info-list">
              <view class="info-item">
                <text class="item-label">身份证号</text>
                <text class="item-value">{{ maskIdCard(elderlyInfo.idCard) }}</text>
              </view>
              <view class="info-item">
                <text class="item-label">出生日期</text>
                <text class="item-value">{{ elderlyInfo.birthDate }}</text>
              </view>
              <view class="info-item">
                <text class="item-label">入住日期</text>
                <text class="item-value">{{ elderlyInfo.checkInDate }}</text>
              </view>
              <view class="info-item">
                <text class="item-label">联系电话</text>
                <text class="item-value">{{ elderlyInfo.phone }}</text>
              </view>
            </view>
          </view>

          <view class="info-section">
            <view class="section-title">家属信息</view>
            <view class="family-list">
              <view class="family-item" v-for="family in elderlyInfo.family" :key="family.id">
                <view class="family-info">
                  <text class="family-name">{{ family.name }}</text>
                  <text class="family-relation">{{ family.relation }}</text>
                </view>
                <view class="family-actions">
                  <view class="action-btn" @click="handleCall(family.phone)">
                    <uni-icons type="phone" size="16" color="#409eff" />
                    <text>呼叫</text>
                  </view>
                </view>
              </view>
            </view>
          </view>

          <view class="info-section">
            <view class="section-title">健康状况</view>
            <text class="health-desc">{{ elderlyInfo.healthStatus || '暂无记录' }}</text>
          </view>
        </view>

        <!-- 护理记录 -->
        <view v-if="activeTab === 'care'" class="tab-content">
          <view class="care-list">
            <view class="care-item" v-for="record in careRecords" :key="record.id" @click="handleCareDetail(record)">
              <view class="care-header">
                <view class="care-type-tag">{{ record.type }}</view>
                <text class="care-time">{{ formatTime(record.time) }}</text>
              </view>
              <text class="care-content">{{ record.content }}</text>
              <view class="care-footer">
                <text class="care-nurse">{{ record.nurseName }}</text>
                <view class="care-rating" v-if="record.evaluation">
                  <uni-rate :value="record.evaluation" size="12" readonly />
                </view>
              </view>
            </view>
          </view>

          <view v-if="careRecords.length === 0" class="empty-state">
            <text>暂无护理记录</text>
          </view>

          <view class="load-more" v-if="hasMore" @click="loadMoreCare">
            <text>加载更多</text>
          </view>
        </view>

        <!-- 健康数据 -->
        <view v-if="activeTab === 'health'" class="tab-content">
          <view class="health-cards">
            <view class="health-card" v-for="item in healthData" :key="item.key">
              <view class="card-header-row">
                <text class="card-label">{{ item.label }}</text>
                <view class="latest-badge">最新</view>
              </view>
              <view class="card-value-row">
                <text class="card-value">{{ item.value }}</text>
                <text class="card-unit">{{ item.unit }}</text>
              </view>
              <text class="card-time">{{ item.time }}</text>
            </view>
          </view>

          <view class="chart-section">
            <view class="section-title">健康趋势</view>
            <view class="chart-placeholder">
              <text class="placeholder-text">健康趋势图表</text>
              <text class="placeholder-desc">（需要引入图表组件）</text>
            </view>
          </view>
        </view>

        <!-- 费用账单 -->
        <view v-if="activeTab === 'bill'" class="tab-content">
          <view class="bill-list">
            <view class="bill-item" v-for="bill in bills" :key="bill.id" @click="handleBillDetail(bill)">
              <view class="bill-header">
                <text class="bill-no">{{ bill.billNo }}</text>
                <view class="bill-status" :class="`status-${bill.status}`">
                  {{ getStatusText(bill.status) }}
                </view>
              </view>
              <view class="bill-body">
                <view class="bill-info-row">
                  <text class="bill-type">{{ bill.billType }}</text>
                  <text class="bill-amount">¥{{ bill.amount.toFixed(2) }}</text>
                </view>
                <text class="bill-date">应付日期：{{ bill.dueDate }}</text>
              </view>
              <view class="bill-actions" v-if="bill.status === 'unpaid'">
                <button class="pay-btn" @click.stop="handlePay(bill)">立即支付</button>
              </view>
            </view>
          </view>

          <view v-if="bills.length === 0" class="empty-state">
            <text>暂无账单记录</text>
          </view>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { formatRelativeTime, maskIdCard } from '@/utils'

// 老人信息
const elderlyInfo = ref<any>({
  id: '',
  name: '',
  avatar: '',
  gender: '',
  age: 0,
  bedNumber: '',
  careLevel: '',
  birthDate: '',
  idCard: '',
  checkInDate: '',
  stayDays: 0,
  careCount: 0,
  healthScore: 0,
  phone: '',
  healthStatus: '',
  family: []
})

// 标签页
const activeTab = ref('info')
const tabs = ref([
  { key: 'info', label: '基本信息' },
  { key: 'care', label: '护理记录' },
  { key: 'health', label: '健康数据' },
  { key: 'bill', label: '费用账单' }
])

// 护理记录
const careRecords = ref<any[]>([])
const hasMore = ref(false)

// 健康数据
const healthData = ref([
  { key: 'bp', label: '血压', value: '128/82', unit: 'mmHg', time: '今日 14:30' },
  { key: 'hr', label: '心率', value: '72', unit: '次/分', time: '今日 14:30' },
  { key: 'bs', label: '血糖', value: '6.3', unit: 'mmol/L', time: '今日 08:00' },
  { key: 'temp', label: '体温', value: '36.5', unit: '℃', time: '今日 14:30' }
])

// 账单
const bills = ref([
  {
    id: 1,
    billNo: 'B202603001',
    billType: '床位费',
    amount: 3500,
    dueDate: '2026-03-10',
    status: 'unpaid'
  },
  {
    id: 2,
    billNo: 'B202603002',
    billType: '护理费',
    amount: 1800,
    dueDate: '2026-03-10',
    status: 'paid'
  }
])

const getCareLevelText = (level: string) => {
  const map: Record<string, string> = {
    level3: '三级护理',
    level2: '二级护理',
    level1: '一级护理',
    special: '特级护理'
  }
  return map[level] || level
}

const formatTime = (timestamp: number) => {
  return formatRelativeTime(timestamp)
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    paid: '已支付',
    unpaid: '未支付',
    overdue: '已逾期'
  }
  return map[status] || status
}

const handleCall = (phone: string) => {
  uni.makePhoneCall({
    phoneNumber: phone
  })
}

const handleCareDetail = (record: any) => {
  console.log('护理记录详情', record)
}

const loadMoreCare = () => {
  console.log('加载更多')
}

const handleBillDetail = (bill: any) => {
  uni.navigateTo({
    url: `/pages/bills/detail?id=${bill.id}`
  })
}

const handlePay = (bill: any) => {
  uni.navigateTo({
    url: `/pages/bills/payment?id=${bill.id}`
  })
}

const loadData = async () => {
  // TODO: 加载老人详情数据
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const options = currentPage.options || {}
  const elderlyId = options.id

  if (!elderlyId) {
    uni.showToast({
      title: '参数错误',
      icon: 'none'
    })
    return
  }

  // 模拟数据
  elderlyInfo.value = {
    id: elderlyId,
    name: '张奶奶',
    avatar: '',
    gender: '女',
    age: 78,
    bedNumber: '3号楼201',
    careLevel: 'level2',
    birthDate: '1946-01-01',
    idCard: '110101194601011234',
    checkInDate: '2023-06-15',
    stayDays: 268,
    careCount: 128,
    healthScore: 85,
    phone: '138****1234',
    healthStatus: '老人整体状况良好，血压稳定，食欲正常，睡眠质量良好。',
    family: [
      { id: 1, name: '张先生', relation: '儿子', phone: '13812345678' },
      { id: 2, name: '李女士', relation: '女儿', phone: '13987654321' }
    ]
  }

  // 模拟护理记录
  careRecords.value = [
    {
      id: 1,
      type: '日常护理',
      content: '测量血压、体温，协助用餐',
      time: Date.now() - 3600000,
      nurseName: '赵护士',
      evaluation: 5
    },
    {
      id: 2,
      type: '康复训练',
      content: '上肢关节活动训练30分钟',
      time: Date.now() - 7200000,
      nurseName: '陈康复师',
      evaluation: 4
    }
  ]
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
.elderly-info-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.elderly-card {
  position: relative;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30rpx;
  overflow: hidden;
}

.card-bg {
  position: absolute;
  top: -100rpx;
  right: -100rpx;
  width: 400rpx;
  height: 400rpx;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
}

.card-content {
  position: relative;
  z-index: 1;
}

.elderly-avatar-wrapper {
  position: relative;
  width: 160rpx;
  height: 160rpx;
  margin: 0 auto 20rpx;
}

.elderly-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 6rpx solid rgba(255, 255, 255, 0.3);
}

.elderly-avatar-placeholder {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 6rpx solid rgba(255, 255, 255, 0.3);
}

.avatar-text {
  font-size: 60rpx;
  color: #fff;
  font-weight: 600;
}

.care-level-tag {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  padding: 8rpx 20rpx;
  border-radius: 20rpx;
  font-size: 22rpx;
  color: #fff;

  &.level-level3 {
    background: rgba(144, 147, 153, 0.9);
  }

  &.level-level2 {
    background: rgba(103, 194, 58, 0.9);
  }

  &.level-level1 {
    background: rgba(230, 162, 60, 0.9);
  }

  &.level-special {
    background: rgba(245, 108, 108, 0.9);
  }
}

.elderly-base-info {
  text-align: center;
  margin-bottom: 30rpx;
}

.elderly-name {
  display: block;
  font-size: 40rpx;
  font-weight: 600;
  color: #fff;
  margin-bottom: 16rpx;
}

.info-tags {
  display: flex;
  justify-content: center;
  gap: 12rpx;
}

.info-tag {
  padding: 8rpx 20rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 20rpx;
}

.tag-text {
  font-size: 24rpx;
  color: #fff;
}

.elderly-stats {
  display: flex;
  justify-content: space-around;
  padding: 24rpx 0;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 16rpx;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.stat-value {
  font-size: 36rpx;
  font-weight: 600;
  color: #fff;
}

.stat-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
}

.stat-divider {
  width: 1rpx;
  background: rgba(255, 255, 255, 0.3);
}

.tabs-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.tabs-header {
  display: flex;
  background: #fff;
  padding: 0 30rpx;
  border-bottom: 1rpx solid #ebeef5;
}

.tab-item {
  position: relative;
  padding: 24rpx 0;
  margin-right: 40rpx;
}

.tab-text {
  font-size: 28rpx;
  color: #606266;

  .active & {
    color: #409eff;
    font-weight: 500;
  }
}

.tab-indicator {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 3rpx;
  background: #409eff;
  border-radius: 2rpx;
}

.tabs-content {
  flex: 1;
  padding: 30rpx;
}

.tab-content {
  .info-section {
    background: #fff;
    border-radius: 16rpx;
    padding: 24rpx;
    margin-bottom: 20rpx;
  }

  .section-title {
    font-size: 28rpx;
    font-weight: 600;
    color: #303133;
    margin-bottom: 20rpx;
  }

  .info-list {
    .info-item {
      display: flex;
      justify-content: space-between;
      padding: 16rpx 0;
      border-bottom: 1rpx solid #f5f7fa;

      &:last-child {
        border-bottom: none;
      }
    }

    .item-label {
      font-size: 28rpx;
      color: #909399;
    }

    .item-value {
      font-size: 28rpx;
      color: #303133;
    }
  }

  .family-list {
    .family-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20rpx 0;
      border-bottom: 1rpx solid #f5f7fa;

      &:last-child {
        border-bottom: none;
      }
    }

    .family-info {
      display: flex;
      flex-direction: column;
      gap: 8rpx;
    }

    .family-name {
      font-size: 28rpx;
      font-weight: 500;
      color: #303133;
    }

    .family-relation {
      font-size: 24rpx;
      color: #909399;
    }

    .action-btn {
      display: flex;
      align-items: center;
      gap: 8rpx;
      padding: 12rpx 24rpx;
      background: #ecf5ff;
      border-radius: 8rpx;
      font-size: 24rpx;
      color: #409eff;
    }
  }

  .health-desc {
    font-size: 28rpx;
    color: #606266;
    line-height: 1.8;
  }
}

.care-list {
  .care-item {
    background: #fff;
    border-radius: 16rpx;
    padding: 24rpx;
    margin-bottom: 20rpx;
  }

  .care-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16rpx;
  }

  .care-type-tag {
    padding: 8rpx 16rpx;
    background: #ecf5ff;
    color: #409eff;
    font-size: 24rpx;
    border-radius: 8rpx;
  }

  .care-time {
    font-size: 24rpx;
    color: #909399;
  }

  .care-content {
    display: block;
    font-size: 28rpx;
    color: #303133;
    line-height: 1.6;
    margin-bottom: 16rpx;
  }

  .care-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .care-nurse {
    font-size: 24rpx;
    color: #909399;
  }
}

.health-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
  margin-bottom: 30rpx;
}

.health-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.card-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.card-label {
  font-size: 24rpx;
  color: #909399;
}

.latest-badge {
  padding: 4rpx 12rpx;
  background: #ecf5ff;
  color: #409eff;
  font-size: 20rpx;
  border-radius: 8rpx;
}

.card-value-row {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
  margin-bottom: 12rpx;
}

.card-value {
  font-size: 44rpx;
  font-weight: 600;
  color: #303133;
}

.card-unit {
  font-size: 24rpx;
  color: #909399;
}

.card-time {
  font-size: 22rpx;
  color: #c0c4cc;
}

.chart-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.chart-placeholder {
  padding: 60rpx 0;
  text-align: center;
  background: #f5f7fa;
  border-radius: 12rpx;
}

.placeholder-text {
  display: block;
  font-size: 28rpx;
  color: #909399;
  margin-bottom: 8rpx;
}

.placeholder-desc {
  display: block;
  font-size: 24rpx;
  color: #c0c4cc;
}

.bill-list {
  .bill-item {
    background: #fff;
    border-radius: 16rpx;
    padding: 24rpx;
    margin-bottom: 20rpx;
  }

  .bill-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16rpx;
  }

  .bill-no {
    font-size: 26rpx;
    font-weight: 500;
    color: #303133;
  }

  .bill-status {
    padding: 8rpx 16rpx;
    border-radius: 8rpx;
    font-size: 22rpx;

    &.status-paid {
      background: #f0f9ff;
      color: #67c23a;
    }

    &.status-unpaid {
      background: #fdf6ec;
      color: #e6a23c;
    }

    &.status-overdue {
      background: #fef0f0;
      color: #f56c6c;
    }
  }

  .bill-body {
    margin-bottom: 16rpx;
  }

  .bill-info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12rpx;
  }

  .bill-type {
    font-size: 24rpx;
    color: #606266;
    padding: 4rpx 12rpx;
    background: #f5f7fa;
    border-radius: 8rpx;
  }

  .bill-amount {
    font-size: 32rpx;
    font-weight: 600;
    color: #f56c6c;
  }

  .bill-date {
    font-size: 24rpx;
    color: #909399;
  }

  .bill-actions {
    display: flex;
    justify-content: flex-end;
  }

  .pay-btn {
    padding: 16rpx 40rpx;
    background: #409eff;
    color: #fff;
    font-size: 28rpx;
    border: none;
    border-radius: 8rpx;
  }
}

.empty-state {
  padding: 80rpx 0;
  text-align: center;
  font-size: 28rpx;
  color: #909399;
}

.load-more {
  padding: 30rpx;
  text-align: center;
  font-size: 28rpx;
  color: #409eff;
}
</style>
