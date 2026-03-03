<template>
  <view class="health-page">
    <!-- 头部 -->
    <view class="page-header">
      <text class="header-title">健康数据</text>
      <view class="header-actions">
        <view class="action-btn" @click="handleAddRecord">
          <uni-icons type="plus" size="18" color="#409eff" />
          <text>记录</text>
        </view>
      </view>
    </view>

    <!-- 最新数据 -->
    <view class="latest-section">
      <view class="section-title">今日数据</view>
      <view class="health-cards">
        <view class="health-card" v-for="item in healthItems" :key="item.key">
          <view class="card-icon" :class="`icon-${item.level}`">
            <uni-icons :type="item.icon" size="24" color="#fff" />
          </view>
          <view class="card-content">
            <text class="card-label">{{ item.label }}</text>
            <view class="card-value-row">
              <text class="card-value">{{ item.value }}</text>
              <text class="card-unit">{{ item.unit }}</text>
            </view>
            <text class="card-time">{{ item.time }}</text>
          </view>
          <view class="card-trend" :class="item.trend > 0 ? 'up' : 'down'">
            <uni-icons :type="item.trend > 0 ? 'arrow-up' : 'arrow-down'" size="14" />
            <text>{{ Math.abs(item.trend) }}%</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 健康趋势图 -->
    <view class="chart-section">
      <view class="section-header">
        <text class="section-title">健康趋势</text>
        <view class="type-selector">
          <view
            v-for="type in chartTypes"
            :key="type.value"
            class="type-item"
            :class="{ active: selectedType === type.value }"
            @click="selectedType = type.value"
          >
            <text class="type-text">{{ type.label }}</text>
          </view>
        </view>
      </view>
      <view class="chart-placeholder">
        <text class="placeholder-text">健康趋势图表</text>
        <text class="placeholder-desc">（需要引入图表组件）</text>
      </view>
    </view>

    <!-- 历史记录 -->
    <view class="history-section">
      <view class="section-header">
        <text class="section-title">历史记录</text>
        <view class="filter-btn" @click="handleFilter">
          <uni-icons type="settings" size="14" />
          <text>筛选</text>
        </view>
      </view>

      <view class="date-list">
        <view class="date-item" v-for="record in historyRecords" :key="record.id">
          <view class="date-header">
            <text class="date-text">{{ record.date }}</text>
            <uni-icons type="right" size="16" color="#c0c4cc" />
          </view>
          <view class="date-data">
            <view class="data-item" v-for="item in record.data" :key="item.key">
              <text class="data-label">{{ item.label }}</text>
              <text class="data-value" :class="`value-${item.level}`">{{ item.value }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const selectedType = ref('bloodPressure')

const chartTypes = ref([
  { label: '血压', value: 'bloodPressure' },
  { label: '血糖', value: 'bloodSugar' },
  { label: '心率', value: 'heartRate' },
  { label: '体温', value: 'temperature' }
])

const healthItems = ref([
  {
    key: 'bp',
    label: '血压',
    value: '128/82',
    unit: 'mmHg',
    icon: 'pulse',
    level: 'normal',
    trend: -2,
    time: '14:30'
  },
  {
    key: 'hr',
    label: '心率',
    value: '72',
    unit: '次/分',
    icon: 'heart',
    level: 'normal',
    trend: 0,
    time: '14:30'
  },
  {
    key: 'bs',
    label: '血糖',
    value: '6.3',
    unit: 'mmol/L',
    icon: 'waterdrop',
    level: 'normal',
    trend: 5,
    time: '08:00'
  },
  {
    key: 'temp',
    label: '体温',
    value: '36.5',
    unit: '℃',
    icon: 'thermometer',
    level: 'normal',
    trend: 0,
    time: '14:30'
  },
  {
    key: 'weight',
    label: '体重',
    value: '58',
    unit: 'kg',
    icon: 'person',
    level: 'normal',
    trend: -1,
    time: '06:00'
  },
  {
    key: 'spo2',
    label: '血氧',
    value: '97',
    unit: '%',
    icon: 'loop',
    level: 'normal',
    trend: 0,
    time: '14:30'
  }
])

const historyRecords = ref([
  {
    id: 1,
    date: '2026-03-03',
    data: [
      { key: 'bp', label: '血压', value: '128/82', level: 'normal' },
      { key: 'bs', label: '血糖', value: '6.3', level: 'normal' }
    ]
  },
  {
    id: 2,
    date: '2026-03-02',
    data: [
      { key: 'bp', label: '血压', value: '130/85', level: 'normal' },
      { key: 'bs', label: '血糖', value: '6.5', level: 'normal' }
    ]
  },
  {
    id: 3,
    date: '2026-03-01',
    data: [
      { key: 'bp', label: '血压', value: '126/80', level: 'normal' },
      { key: 'bs', label: '血糖', value: '6.2', level: 'normal' }
    ]
  }
])

const handleAddRecord = () => {
  uni.navigateTo({
    url: '/pages/health/add'
  })
}

const handleFilter = () => {
  uni.showActionSheet({
    itemList: ['最近7天', '最近30天', '最近3个月'],
    success: (res) => {
      console.log('筛选时间范围', res.tapIndex)
    }
  })
}

onMounted(() => {
  // 加载数据
})
</script>

<style lang="scss" scoped>
.health-page {
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

.action-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 12rpx 24rpx;
  background: #ecf5ff;
  color: #409eff;
  border-radius: 8rpx;
  font-size: 26rpx;
}

.latest-section {
  padding: 30rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20rpx;
}

.health-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.health-card {
  position: relative;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  display: flex;
  align-items: center;
  gap: 16rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.card-icon {
  width: 64rpx;
  height: 64rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;

  &.icon-normal {
    background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
  }

  &.icon-warning {
    background: linear-gradient(135deg, #fccb90 0%, #d57eeb 100%);
  }

  &.icon-danger {
    background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);
  }
}

.card-content {
  flex: 1;
}

.card-label {
  display: block;
  font-size: 22rpx;
  color: #909399;
  margin-bottom: 4rpx;
}

.card-value-row {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
  margin-bottom: 4rpx;
}

.card-value {
  font-size: 32rpx;
  font-weight: 600;
  color: #303133;
}

.card-unit {
  font-size: 22rpx;
  color: #909399;
}

.card-time {
  display: block;
  font-size: 20rpx;
  color: #c0c4cc;
}

.card-trend {
  display: flex;
  align-items: center;
  gap: 4rpx;
  padding: 8rpx;
  border-radius: 8rpx;
  font-size: 20rpx;

  &.up {
    background: #fef0f0;
    color: #f56c6c;
  }

  &.down {
    background: #f0f9ff;
    color: #67c23a;
  }
}

.chart-section {
  padding: 0 30rpx 30rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.type-selector {
  display: flex;
  gap: 12rpx;
}

.type-item {
  padding: 8rpx 20rpx;
  background: #fff;
  border-radius: 20rpx;
  border: 1rpx solid #dcdfe6;

  &.active {
    background: #409eff;
    border-color: #409eff;
  }

  .type-text {
    font-size: 24rpx;

    .active & {
      color: #fff;
    }
  }
}

.chart-placeholder {
  background: #fff;
  border-radius: 16rpx;
  padding: 80rpx 0;
  text-align: center;
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

.history-section {
  padding: 0 30rpx 30rpx;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 26rpx;
  color: #606266;
}

.date-list {
  .date-item {
    background: #fff;
    border-radius: 16rpx;
    padding: 24rpx;
    margin-bottom: 16rpx;
  }
}

.date-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid #f5f7fa;
}

.date-text {
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;
}

.date-data {
  display: flex;
  flex-wrap: wrap;
  gap: 24rpx;
}

.data-item {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.data-label {
  font-size: 24rpx;
  color: #909399;
}

.data-value {
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;

  &.value-warning {
    color: #e6a23c;
  }

  &.value-danger {
    color: #f56c6c;
  }
}
</style>
