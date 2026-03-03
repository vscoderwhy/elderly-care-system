<template>
  <view class="appointment-page">
    <!-- 顶部状态栏 -->
    <view class="status-bar">
      <view class="status-back" @click="handleBack">
        <uni-icons type="back" size="20" color="#606266" />
      </view>
      <text class="status-title">{{ isEdit ? '编辑预约' : '探视预约' }}</text>
      <view class="status-placeholder"></view>
    </view>

    <scroll-view class="appointment-content" scroll-y>
      <!-- 预约表单 -->
      <view class="form-section">
        <view class="form-group">
          <view class="form-label required">预约老人</view>
          <picker :value="elderlyIndex" :range="elderlyList" range-key="label" @change="handleElderlyChange">
            <view class="picker-value" :class="{ placeholder: !form.elderlyId }">
              {{ selectedElderly || '请选择老人' }}
              <uni-icons type="arrowright" size="14" color="#909399" />
            </view>
          </picker>
        </view>

        <view class="form-group">
          <view class="form-label required">访客姓名</view>
          <input
            class="form-input"
            v-model="form.visitorName"
            placeholder="请输入访客姓名"
            placeholder-class="input-placeholder"
          />
        </view>

        <view class="form-group">
          <view class="form-label required">联系电话</view>
          <input
            class="form-input"
            v-model="form.visitorPhone"
            type="number"
            maxlength="11"
            placeholder="请输入联系电话"
            placeholder-class="input-placeholder"
          />
        </view>

        <view class="form-group">
          <view class="form-label required">与老人关系</view>
          <picker :value="relationIndex" :range="relations" range-key="label" @change="handleRelationChange">
            <view class="picker-value" :class="{ placeholder: !form.relationship }">
              {{ form.relationship || '请选择关系' }}
              <uni-icons type="arrowright" size="14" color="#909399" />
            </view>
          </picker>
        </view>

        <view class="form-group">
          <view class="form-label required">探访类型</view>
          <view class="radio-group">
            <view
              class="radio-item"
              :class="{ active: form.visitType === '现场探访' }"
              @click="form.visitType = '现场探访'"
            >
              <view class="radio-icon">
                <uni-icons v-if="form.visitType === '现场探访'" type="checkmarkempty" size="16" color="#409eff" />
              </view>
              <text>现场探访</text>
            </view>
            <view
              class="radio-item"
              :class="{ active: form.visitType === '视频探访' }"
              @click="form.visitType = '视频探访'"
            >
              <view class="radio-icon">
                <uni-icons v-if="form.visitType === '视频探访'" type="checkmarkempty" size="16" color="#409eff" />
              </view>
              <text>视频探访</text>
            </view>
          </view>
        </view>

        <view class="form-group">
          <view class="form-label required">预约日期</view>
          <picker mode="date" :value="form.visitDate" :start="minDate" :end="maxDate" @change="handleDateChange">
            <view class="picker-value" :class="{ placeholder: !form.visitDate }">
              {{ form.visitDate || '请选择日期' }}
              <uni-icons type="arrowright" size="14" color="#909399" />
            </view>
          </picker>
        </view>

        <view class="form-group">
          <view class="form-label required">预约时间</view>
          <picker :value="timeIndex" :range="timeSlots" range-key="label" @change="handleTimeChange">
            <view class="picker-value" :class="{ placeholder: !form.visitTime }">
              {{ form.visitTime || '请选择时间' }}
              <uni-icons type="arrowright" size="14" color="#909399" />
            </view>
          </picker>
        </view>

        <view class="form-group">
          <view class="form-label required">探访人数</view>
          <view class="counter-wrapper">
            <view class="counter-btn" @click="handleCountChange(-1)">
              <uni-icons type="minus" size="16" color="#606266" />
            </view>
            <text class="counter-value">{{ form.visitorCount }}</text>
            <view class="counter-btn" @click="handleCountChange(1)">
              <uni-icons type="plus" size="16" color="#606266" />
            </view>
          </view>
        </view>

        <view class="form-group">
          <view class="form-label">备注信息</view>
          <textarea
            class="form-textarea"
            v-model="form.remark"
            placeholder="请输入备注信息（选填）"
            placeholder-class="input-placeholder"
            maxlength="200"
          />
          <view class="textarea-counter">{{ form.remark.length }}/200</view>
        </view>
      </view>

      <!-- 预约须知 -->
      <view class="notice-section">
        <view class="notice-title">
          <uni-icons type="info" size="16" color="#e6a23c" />
          <text>预约须知</text>
        </view>
        <view class="notice-content">
          <text>1. 探访时间为：09:00-11:00，14:00-17:00</text>
          <text>2. 现场探访需提前1天预约</text>
          <text>3. 每次探访时间不超过1小时</text>
          <text>4. 探访人数不超过3人</text>
          <text>5. 如需取消，请提前2小时联系</text>
        </view>
      </view>
    </scroll-view>

    <!-- 底部提交栏 -->
    <view class="submit-bar">
      <button class="submit-btn" @click="handleSubmit" :disabled="submitting">
        {{ submitting ? '提交中...' : '提交预约' }}
      </button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

interface Elderly {
  id: number
  label: string
  name: string
  bedNumber: string
}

const isEdit = ref(false)
const elderlyIndex = ref(0)
const relationIndex = ref(0)
const timeIndex = ref(0)
const submitting = ref(false)

const elderlyList = ref<Elderly[]>([
  { id: 1, label: '张奶奶 - 3号楼201', name: '张奶奶', bedNumber: '3号楼201' },
  { id: 2, label: '李爷爷 - 2号楼105', name: '李爷爷', bedNumber: '2号楼105' },
  { id: 3, label: '王奶奶 - 3号楼202', name: '王奶奶', bedNumber: '3号楼202' }
])

const relations = ref([
  { value: '子女', label: '子女' },
  { value: '配偶', label: '配偶' },
  { value: '孙辈', label: '孙辈' },
  { value: '其他亲属', label: '其他亲属' },
  { value: '朋友', label: '朋友' },
  { value: '其他', label: '其他' }
])

const timeSlots = ref([
  { value: '09:00', label: '09:00-10:00' },
  { value: '10:00', label: '10:00-11:00' },
  { value: '14:00', label: '14:00-15:00' },
  { value: '15:00', label: '15:00-16:00' },
  { value: '16:00', label: '16:00-17:00' }
])

const form = ref({
  elderlyId: '',
  visitorName: '',
  visitorPhone: '',
  relationship: '',
  visitType: '现场探访',
  visitDate: '',
  visitTime: '',
  visitorCount: 1,
  remark: ''
})

const minDate = computed(() => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  return tomorrow.toISOString().split('T')[0]
})

const maxDate = computed(() => {
  const maxDay = new Date()
  maxDay.setDate(maxDay.getDate() + 30)
  return maxDay.toISOString().split('T')[0]
})

const selectedElderly = computed(() => {
  if (!form.value.elderlyId) return ''
  const elderly = elderlyList.value.find(e => e.id === parseInt(form.value.elderlyId))
  return elderly?.label || ''
})

const handleBack = () => {
  uni.navigateBack()
}

const handleElderlyChange = (e: any) => {
  elderlyIndex.value = e.detail.value
  form.value.elderlyId = String(elderlyList.value[e.detail.value].id)
}

const handleRelationChange = (e: any) => {
  relationIndex.value = e.detail.value
  form.value.relationship = relations.value[e.detail.value].value
}

const handleDateChange = (e: any) => {
  form.value.visitDate = e.detail.value
}

const handleTimeChange = (e: any) => {
  timeIndex.value = e.detail.value
  form.value.visitTime = timeSlots.value[e.detail.value].value
}

const handleCountChange = (delta: number) => {
  const newValue = form.value.visitorCount + delta
  if (newValue >= 1 && newValue <= 10) {
    form.value.visitorCount = newValue
  }
}

const validateForm = () => {
  if (!form.value.elderlyId) {
    uni.showToast({ title: '请选择预约老人', icon: 'none' })
    return false
  }
  if (!form.value.visitorName.trim()) {
    uni.showToast({ title: '请输入访客姓名', icon: 'none' })
    return false
  }
  if (!/^1[3-9]\d{9}$/.test(form.value.visitorPhone)) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return false
  }
  if (!form.value.relationship) {
    uni.showToast({ title: '请选择与老人关系', icon: 'none' })
    return false
  }
  if (!form.value.visitDate) {
    uni.showToast({ title: '请选择预约日期', icon: 'none' })
    return false
  }
  if (!form.value.visitTime) {
    uni.showToast({ title: '请选择预约时间', icon: 'none' })
    return false
  }
  return true
}

const handleSubmit = async () => {
  if (!validateForm()) return

  submitting.value = true

  // TODO: 调用API提交预约
  setTimeout(() => {
    submitting.value = false
    uni.showToast({
      title: '预约提交成功',
      icon: 'success'
    })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  }, 1000)
}

onLoad((options: any) => {
  if (options.id) {
    isEdit.value = true
    // TODO: 加载预约详情
  }
})
</script>

<style lang="scss" scoped>
.appointment-page {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 120rpx;
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

.appointment-content {
  padding: 20rpx 30rpx;
}

.form-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.form-group {
  margin-bottom: 30rpx;

  &:last-child {
    margin-bottom: 0;
  }
}

.form-label {
  font-size: 28rpx;
  color: #303133;
  margin-bottom: 16rpx;
  display: flex;
  align-items: center;

  &.required::before {
    content: '*';
    color: #f56c6c;
    margin-right: 4rpx;
  }
}

.form-input {
  width: 100%;
  height: 80rpx;
  padding: 0 24rpx;
  background: #f5f7fa;
  border-radius: 8rpx;
  font-size: 28rpx;
  color: #303133;
}

.input-placeholder {
  color: #909399;
}

.picker-value {
  height: 80rpx;
  padding: 0 24rpx;
  background: #f5f7fa;
  border-radius: 8rpx;
  font-size: 28rpx;
  color: #303133;
  display: flex;
  align-items: center;
  justify-content: space-between;

  &.placeholder {
    color: #909399;
  }
}

.radio-group {
  display: flex;
  gap: 20rpx;
}

.radio-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 20rpx 24rpx;
  background: #f5f7fa;
  border-radius: 8rpx;
  border: 2px solid transparent;
  transition: all 0.3s;

  &.active {
    border-color: #409eff;
    background: #ecf5ff;
  }
}

.radio-icon {
  width: 32rpx;
  height: 32rpx;
  border-radius: 50%;
  border: 2px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: center;

  .radio-item.active & {
    border-color: #409eff;
    background: #409eff;
  }
}

.counter-wrapper {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 16rpx 24rpx;
  background: #f5f7fa;
  border-radius: 8rpx;
  width: fit-content;
}

.counter-btn {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #dcdfe6;
}

.counter-value {
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;
  min-width: 40rpx;
  text-align: center;
}

.form-textarea {
  width: 100%;
  min-height: 160rpx;
  padding: 16rpx 24rpx;
  background: #f5f7fa;
  border-radius: 8rpx;
  font-size: 28rpx;
  color: #303133;
}

.textarea-counter {
  text-align: right;
  font-size: 24rpx;
  color: #909399;
  margin-top: 8rpx;
}

.notice-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
}

.notice-title {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16rpx;
}

.notice-content {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.notice-content text {
  font-size: 26rpx;
  color: #606266;
  line-height: 1.6;
}

.submit-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 30rpx;
  background: #fff;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.submit-btn {
  width: 100%;
  height: 80rpx;
  background: #409eff;
  color: #fff;
  border: none;
  border-radius: 40rpx;
  font-size: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;

  &:disabled {
    background: #c0c4cc;
  }
}
</style>
