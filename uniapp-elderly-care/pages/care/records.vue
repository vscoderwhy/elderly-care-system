<template>
  <view class="care-records-page">
    <!-- 页面头部 -->
    <view class="page-header">
      <text class="header-title">护理记录</text>
    </view>

    <!-- 筛选栏 -->
    <view class="filter-bar">
      <view class="filter-item">
        <uni-data-select
          v-model="filterForm.careType"
          :localdata="careTypes"
          placeholder="护理类型"
          @change="loadRecords"
        />
      </view>
      <view class="filter-item">
        <uni-data-select
          v-model="filterForm.dateRange"
          :localdata="dateRanges"
          placeholder="时间范围"
          @change="loadRecords"
        />
      </view>
    </view>

    <!-- 记录列表 -->
    <scroll-view class="records-list" scroll-y @scrolltolower="loadMore">
      <view
        class="record-item"
        v-for="record in recordList"
        :key="record.id"
        @click="handleRecordDetail(record)"
      >
        <!-- 头部 -->
        <view class="record-header">
          <view class="care-type-tag" :class="`type-${record.typeCode}`">
            {{ record.careType }}
          </view>
          <text class="record-time">{{ formatTime(record.time) }}</text>
        </view>

        <!-- 内容 -->
        <view class="record-content">
          <text class="record-desc">{{ record.content }}</text>
        </view>

        <!-- 照片 -->
        <view v-if="record.images && record.images.length > 0" class="record-images">
          <image
            v-for="(img, index) in record.images"
            :key="index"
            :src="img"
            class="record-image"
            mode="aspectFill"
            @click.stop="handlePreviewImage(record.images, index)"
          />
        </view>

        <!-- 底部信息 -->
        <view class="record-footer">
          <view class="nurse-info">
            <uni-icons type="person" size="14" color="#909399" />
            <text class="nurse-name">{{ record.nurseName }}</text>
          </view>

          <view class="evaluation-wrapper" v-if="record.evaluation">
            <uni-rate :value="record.evaluation" size="12" readonly />
            <text class="eval-text">{{ record.evaluation }}分</text>
          </view>
        </view>
      </view>

      <!-- 加载更多 -->
      <view class="load-more" v-if="hasMore" @click="loadMore">
        <text>加载更多</text>
      </view>

      <!-- 空状态 -->
      <view v-if="!loading && recordList.length === 0" class="empty-state">
        <image src="/static/empty.png" class="empty-image" mode="aspectFit" />
        <text class="empty-text">暂无护理记录</text>
      </view>

      <!-- 加载中 -->
      <view v-if="loading" class="loading-state">
        <text>加载中...</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { formatRelativeTime } from '@/utils'

// 护理类型
const careTypes = ref([
  { value: '', text: '全部类型' },
  { value: 'daily', text: '日常护理' },
  { value: 'rehab', text: '康复训练' },
  { value: 'health', text: '健康监测' },
  { value: 'medical', text: '医疗护理' }
])

// 时间范围
const dateRanges = ref([
  { value: '', text: '全部时间' },
  { value: '7d', text: '最近7天' },
  { value: '30d', text: '最近30天' },
  { value: '90d', text: '最近3个月' }
])

const filterForm = reactive({
  careType: '',
  dateRange: ''
})

const recordList = ref<any[]>([])
const loading = ref(false)
const hasMore = ref(false)
const pagination = reactive({
  page: 1,
  pageSize: 20
})

const getTypeCode = (type: string) => {
  const map: Record<string, string> = {
    '日常护理': 'daily',
    '康复训练': 'rehab',
    '健康监测': 'health',
    '医疗护理': 'medical'
  }
  return map[type] || 'daily'
}

const formatTime = (timestamp: number) => {
  return formatRelativeTime(timestamp)
}

const loadRecords = async () => {
  loading.value = true
  try {
    // TODO: 实际 API 请求
    await new Promise(resolve => setTimeout(resolve, 500))

    // 模拟数据
    recordList.value = [
      {
        id: 1,
        careType: '日常护理',
        typeCode: 'daily',
        content: '测量血压、体温，协助用餐，老人状态良好，配合度佳',
        time: Date.now() - 3600000,
        images: [],
        nurseName: '赵护士',
        evaluation: 5
      },
      {
        id: 2,
        careType: '康复训练',
        typeCode: 'rehab',
        content: '上肢关节活动训练30分钟，包括屈伸、旋转等动作，老人配合度良好',
        time: Date.now() - 7200000,
        images: ['', ''],
        nurseName: '陈康复师',
        evaluation: 4
      },
      {
        id: 3,
        careType: '健康监测',
        typeCode: 'health',
        content: '血压130/85mmHg，心率72次/分，体温36.4℃',
        time: Date.now() - 86400000,
        images: [],
        nurseName: '周护士',
        evaluation: 0
      },
      {
        id: 4,
        careType: '医疗护理',
        typeCode: 'medical',
        content: '伤口换药，清洁消毒，观察伤口愈合情况',
        time: Date.now() - 172800000,
        images: [''],
        nurseName: '李医生',
        evaluation: 0
      }
    ]

    hasMore.value = recordList.value.length >= pagination.pageSize
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  if (!hasMore.value) return
  pagination.page++
  loadRecords()
}

const handleRecordDetail = (record: any) => {
  uni.navigateTo({
    url: `/pages/care/detail?id=${record.id}`
  })
}

const handlePreviewImage = (images: string[], index: number) => {
  uni.previewImage({
    urls: images,
    current: index
  })
}

onMounted(() => {
  loadRecords()
})
</script>

<style lang="scss" scoped>
.care-records-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 30rpx;
  background: #fff;
}

.header-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #303133;
}

.filter-bar {
  display: flex;
  gap: 12rpx;
  padding: 20rpx 30rpx;
  background: #fff;
}

.filter-item {
  flex: 1;
}

.records-list {
  flex: 1;
  padding: 0 30rpx 30rpx;
}

.record-item {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
}

.record-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.care-type-tag {
  padding: 8rpx 20rpx;
  border-radius: 8rpx;
  font-size: 24rpx;

  &.type-daily {
    background: #e1f3e8;
    color: #67c23a;
  }

  &.type-rehab {
    background: #ecf5ff;
    color: #409eff;
  }

  &.type-health {
    background: #fdf6ec;
    color: #e6a23c;
  }

  &.type-medical {
    background: #fef0f0;
    color: #f56c6c;
  }
}

.record-time {
  font-size: 24rpx;
  color: #909399;
}

.record-content {
  margin-bottom: 16rpx;
}

.record-desc {
  font-size: 28rpx;
  color: #606266;
  line-height: 1.6;
}

.record-images {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.record-image {
  width: 160rpx;
  height: 160rpx;
  border-radius: 8rpx;
}

.record-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nurse-info {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 24rpx;
  color: #909399;
}

.nurse-name {
  color: #606266;
}

.evaluation-wrapper {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.eval-text {
  font-size: 22rpx;
  color: #e6a23c;
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

.empty-image {
  width: 200rpx;
  height: 200rpx;
  margin-bottom: 20rpx;
}

.empty-text {
  display: block;
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
