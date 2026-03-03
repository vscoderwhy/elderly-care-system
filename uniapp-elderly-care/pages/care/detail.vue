<template>
  <view class="care-detail-page">
    <!-- 顶部状态栏 -->
    <view class="status-bar">
      <view class="status-back" @click="handleBack">
        <uni-icons type="back" size="20" color="#606266" />
      </view>
      <text class="status-title">护理详情</text>
      <view class="status-action" @click="handleShare">
        <uni-icons type="share" size="20" color="#606266" />
      </view>
    </view>

    <scroll-view class="detail-content" scroll-y>
      <!-- 基本信息 -->
      <view class="info-section">
        <view class="section-header">
          <view class="care-type-tag" :class="`type-${record.typeCode}`">
            {{ record.careType }}
          </view>
          <text class="record-time">{{ formatTime(record.time) }}</text>
        </view>

        <view class="care-content">
          <text class="content-text">{{ record.content }}</text>
        </view>

        <!-- 护理人员信息 -->
        <view class="nurse-card">
          <view class="nurse-avatar">
            <image :src="record.nurseAvatar || '/static/avatar-default.png'" class="avatar-image" mode="aspectFill" />
          </view>
          <view class="nurse-detail">
            <text class="nurse-name">{{ record.nurseName }}</text>
            <text class="nurse-position">{{ record.nursePosition }}</text>
          </view>
          <view class="nurse-action">
            <button class="call-btn" size="mini" @click="handleCall">
              <uni-icons type="phone" size="14" />
              联系
            </button>
          </view>
        </view>
      </view>

      <!-- 护理图片 -->
      <view class="images-section" v-if="record.images && record.images.length > 0">
        <view class="section-title">
          <text>护理照片</text>
          <text class="image-count">{{ record.images.length }}张</text>
        </view>
        <view class="images-grid">
          <view
            class="image-item"
            v-for="(img, index) in record.images"
            :key="index"
            @click="handlePreviewImage(index)"
          >
            <image :src="img || '/static/image-placeholder.png'" class="care-image" mode="aspectFill" />
          </view>
        </view>
      </view>

      <!-- 护理评价 -->
      <view class="evaluation-section" v-if="record.evaluation > 0">
        <view class="section-title">服务评价</view>
        <view class="evaluation-card">
          <view class="evaluation-stars">
            <uni-rate :value="record.evaluation" size="18" readonly />
            <text class="evaluation-score">{{ record.evaluation }}分</text>
          </view>
          <view class="evaluation-tags" v-if="record.tags && record.tags.length > 0">
            <view
              class="tag-item"
              v-for="(tag, index) in record.tags"
              :key="index"
            >
              {{ tag }}
            </view>
          </view>
          <view class="evaluation-comment" v-if="record.comment">
            <text>{{ record.comment }}</text>
          </view>
        </view>
      </view>

      <!-- 健康数据 -->
      <view class="health-section" v-if="record.healthData">
        <view class="section-title">相关健康数据</view>
        <view class="health-grid">
          <view class="health-item" v-for="(item, key) in record.healthData" :key="key">
            <text class="health-label">{{ getHealthLabel(key) }}</text>
            <text class="health-value">{{ item }}</text>
          </view>
        </view>
      </view>

      <!-- 备注 -->
      <view class="remark-section" v-if="record.remark">
        <view class="section-title">备注</view>
        <view class="remark-content">
          <text>{{ record.remark }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="action-bar" v-if="canEvaluate">
      <button class="action-btn secondary-btn" @click="handleContact">联系护工</button>
      <button class="action-btn primary-btn" @click="handleEvaluate">服务评价</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { formatRelativeTime } from '@/utils'

interface HealthData {
  bloodPressure?: string
  heartRate?: string
  temperature?: string
  bloodSugar?: string
  weight?: string
}

interface CareRecord {
  id: number
  careType: string
  typeCode: string
  content: string
  time: number
  images: string[]
  nurseName: string
  nursePosition: string
  nurseAvatar: string
  evaluation: number
  tags?: string[]
  comment?: string
  healthData?: HealthData
  remark?: string
}

const record = ref<CareRecord>({
  id: 0,
  careType: '',
  typeCode: '',
  content: '',
  time: Date.now(),
  images: [],
  nurseName: '',
  nursePosition: '',
  nurseAvatar: '',
  evaluation: 0
})

const canEvaluate = computed(() => {
  return record.value.evaluation === 0 && Date.now() - record.value.time < 7 * 24 * 60 * 60 * 1000
})

const formatTime = (timestamp: number) => {
  return formatRelativeTime(timestamp)
}

const getHealthLabel = (key: string) => {
  const labels: Record<string, string> = {
    bloodPressure: '血压',
    heartRate: '心率',
    temperature: '体温',
    bloodSugar: '血糖',
    weight: '体重'
  }
  return labels[key] || key
}

const handleBack = () => {
  uni.navigateBack()
}

const handleShare = () => {
  uni.share({
    provider: 'weixin',
    type: 0,
    title: '护理记录分享',
    success: () => {
      uni.showToast({ title: '分享成功', icon: 'success' })
    }
  })
}

const handleCall = () => {
  uni.makePhoneCall({
    phoneNumber: '13800138000'
  })
}

const handleContact = () => {
  uni.showModal({
    title: '联系护工',
    content: `是否联系 ${record.value.nurseName}？`,
    success: (res) => {
      if (res.confirm) {
        handleCall()
      }
    }
  })
}

const handleEvaluate = () => {
  uni.navigateTo({
    url: `/pages/care/evaluate?id=${record.value.id}`
  })
}

const handlePreviewImage = (index: number) => {
  uni.previewImage({
    urls: record.value.images,
    current: index
  })
}

onLoad((options: any) => {
  const id = options.id
  // TODO: 根据ID加载护理记录详情
  // 模拟数据
  record.value = {
    id: parseInt(id) || 1,
    careType: '日常护理',
    typeCode: 'daily',
    content: '协助老人起床洗漱，测量血压、体温、心率，协助用餐，老人状态良好，配合度佳。餐后散步30分钟，情绪稳定。',
    time: Date.now() - 3600000,
    images: [
      '/static/care1.jpg',
      '/static/care2.jpg',
      '/static/care3.jpg'
    ],
    nurseName: '赵护士',
    nursePosition: '高级护理员',
    nurseAvatar: '/static/nurse-avatar.jpg',
    evaluation: 5,
    tags: ['服务热情', '专业细致', '耐心体贴'],
    comment: '赵护士非常有耐心，照顾得很周到，老人很满意。',
    healthData: {
      bloodPressure: '130/85mmHg',
      heartRate: '72次/分',
      temperature: '36.4℃',
      bloodSugar: '6.8mmol/L'
    },
    remark: '老人今日情绪良好，无特殊异常'
  }
})
</script>

<style lang="scss" scoped>
.care-detail-page {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 100rpx;
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
.status-action {
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

.detail-content {
  padding: 20rpx 30rpx;
}

.info-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.care-type-tag {
  padding: 12rpx 24rpx;
  border-radius: 8rpx;
  font-size: 26rpx;
  font-weight: 500;

  &.type-daily { background: #e1f3e8; color: #67c23a; }
  &.type-rehab { background: #ecf5ff; color: #409eff; }
  &.type-health { background: #fdf6ec; color: #e6a23c; }
  &.type-medical { background: #fef0f0; color: #f56c6c; }
}

.record-time {
  font-size: 24rpx;
  color: #909399;
}

.care-content {
  margin-bottom: 30rpx;
}

.content-text {
  font-size: 28rpx;
  color: #606266;
  line-height: 1.8;
}

.nurse-card {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx;
  background: #f5f7fa;
  border-radius: 12rpx;
}

.nurse-avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  overflow: hidden;
}

.avatar-image {
  width: 100%;
  height: 100%;
}

.nurse-detail {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.nurse-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;
}

.nurse-position {
  font-size: 24rpx;
  color: #909399;
}

.call-btn {
  padding: 12rpx 24rpx;
  background: #fff;
  color: #409eff;
  border: 1px solid #409eff;
  border-radius: 8rpx;
  font-size: 24rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.images-section,
.evaluation-section,
.health-section,
.remark-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20rpx;
}

.image-count {
  font-size: 24rpx;
  color: #909399;
  margin-left: 12rpx;
}

.images-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12rpx;
}

.image-item {
  aspect-ratio: 1;
  border-radius: 8rpx;
  overflow: hidden;
}

.care-image {
  width: 100%;
  height: 100%;
}

.evaluation-card {
  padding: 24rpx;
  background: #fef9f3;
  border-radius: 12rpx;
}

.evaluation-stars {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.evaluation-score {
  font-size: 32rpx;
  font-weight: 600;
  color: #e6a23c;
}

.evaluation-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.tag-item {
  padding: 8rpx 20rpx;
  background: #fff;
  color: #e6a23c;
  border: 1px solid #e6a23c;
  border-radius: 20rpx;
  font-size: 24rpx;
}

.evaluation-comment {
  padding: 16rpx;
  background: #fff;
  border-radius: 8rpx;
  font-size: 26rpx;
  color: #606266;
  line-height: 1.6;
}

.health-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.health-item {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
  padding: 20rpx;
  background: #f5f7fa;
  border-radius: 8rpx;
}

.health-label {
  font-size: 24rpx;
  color: #909399;
}

.health-value {
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;
}

.remark-content {
  padding: 24rpx;
  background: #f5f7fa;
  border-radius: 12rpx;
  font-size: 26rpx;
  color: #606266;
  line-height: 1.6;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: 20rpx;
  padding: 20rpx 30rpx;
  background: #fff;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.action-btn {
  flex: 1;
  height: 80rpx;
  border-radius: 40rpx;
  font-size: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
}

.secondary-btn {
  background: #f5f7fa;
  color: #606266;
}

.primary-btn {
  background: #409eff;
  color: #fff;
}
</style>
