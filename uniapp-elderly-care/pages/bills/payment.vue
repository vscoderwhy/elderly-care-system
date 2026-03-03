<template>
  <view class="payment-page">
    <!-- 顶部状态栏 -->
    <view class="status-bar">
      <view class="status-back" @click="handleBack">
        <uni-icons type="back" size="20" color="#606266" />
      </view>
      <text class="status-title">账单支付</text>
      <view class="status-placeholder"></view>
    </view>

    <scroll-view class="payment-content" scroll-y>
      <!-- 账单信息 -->
      <view class="bill-section">
        <view class="section-title">账单信息</view>
        <view class="bill-card">
          <view class="bill-header">
            <text class="bill-no">{{ bill.billNo }}</text>
            <view class="bill-status" :class="`status-${bill.status}`">
              {{ getStatusText(bill.status) }}
            </view>
          </view>
          <view class="bill-amount">
            <text class="amount-label">应付金额</text>
            <text class="amount-value">¥{{ bill.amount.toFixed(2) }}</text>
          </view>
          <view class="bill-details">
            <view class="detail-row">
              <text class="detail-label">账单类型</text>
              <text class="detail-value">{{ bill.billType }}</text>
            </view>
            <view class="detail-row">
              <text class="detail-label">费用期间</text>
              <text class="detail-value">{{ bill.period }}</text>
            </view>
            <view class="detail-row">
              <text class="detail-label">账单日期</text>
              <text class="detail-value">{{ bill.billDate }}</text>
            </view>
            <view class="detail-row">
              <text class="detail-label">应付日期</text>
              <text class="detail-value" :class="{ 'text-danger': isOverdue(bill) }">
                {{ bill.dueDate }}
              </text>
            </view>
          </view>
        </view>
      </view>

      <!-- 费用明细 -->
      <view class="detail-section">
        <view class="section-title">费用明细</view>
        <view class="detail-list">
          <view class="detail-item" v-for="(item, index) in bill.details" :key="index">
            <view class="item-info">
              <text class="item-name">{{ item.name }}</text>
              <text class="item-desc" v-if="item.description">{{ item.description }}</text>
            </view>
            <text class="item-amount">¥{{ item.amount.toFixed(2) }}</text>
          </view>
          <view class="detail-total">
            <text class="total-label">合计</text>
            <text class="total-value">¥{{ bill.amount.toFixed(2) }}</text>
          </view>
        </view>
      </view>

      <!-- 支付方式 -->
      <view class="payment-section">
        <view class="section-title">支付方式</view>
        <view class="payment-methods">
          <view
            class="method-item"
            :class="{ active: selectedMethod === 'wechat' }"
            @click="selectMethod('wechat')"
          >
            <view class="method-icon">
              <text class="icon-wechat">💚</text>
            </view>
            <view class="method-info">
              <text class="method-name">微信支付</text>
            </view>
            <view class="method-check" v-if="selectedMethod === 'wechat'">
              <uni-icons type="checkmarkempty" size="20" color="#07c160" />
            </view>
          </view>

          <view
            class="method-item"
            :class="{ active: selectedMethod === 'alipay' }"
            @click="selectMethod('alipay')"
          >
            <view class="method-icon">
              <text class="icon-alipay">💙</text>
            </view>
            <view class="method-info">
              <text class="method-name">支付宝</text>
            </view>
            <view class="method-check" v-if="selectedMethod === 'alipay'">
              <uni-icons type="checkmarkempty" size="20" color="#1677ff" />
            </view>
          </view>

          <view
            class="method-item"
            :class="{ active: selectedMethod === 'bank' }"
            @click="selectMethod('bank')"
          >
            <view class="method-icon">
              <text class="icon-bank">🏦</text>
            </view>
            <view class="method-info">
              <text class="method-name">银行转账</text>
            </view>
            <view class="method-check" v-if="selectedMethod === 'bank'">
              <uni-icons type="checkmarkempty" size="20" color="#409eff" />
            </view>
          </view>
        </view>
      </view>

      <!-- 支付说明 -->
      <view class="notice-section">
        <view class="notice-title">
          <uni-icons type="info" size="16" color="#e6a23c" />
          <text>支付说明</text>
        </view>
        <view class="notice-content">
          <text>1. 请在应付日期前完成支付，逾期将产生滞纳金</text>
          <text>2. 支付成功后，系统将自动更新账单状态</text>
          <text>3. 如有疑问，请联系客服：400-123-4567</text>
        </view>
      </view>
    </scroll-view>

    <!-- 底部支付栏 -->
    <view class="payment-bar">
      <view class="payment-total">
        <text class="total-label">应付金额：</text>
        <text class="total-amount">¥{{ bill.amount.toFixed(2) }}</text>
      </view>
      <button class="pay-btn" @click="handlePay" :disabled="!selectedMethod || paying">
        {{ paying ? '支付中...' : '立即支付' }}
      </button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

interface BillDetail {
  name: string
  description?: string
  amount: number
}

interface Bill {
  id: number
  billNo: string
  billType: string
  period: string
  billDate: string
  dueDate: string
  amount: number
  status: string
  details: BillDetail[]
}

const bill = ref<Bill>({
  id: 0,
  billNo: '',
  billType: '',
  period: '',
  billDate: '',
  dueDate: '',
  amount: 0,
  status: 'unpaid',
  details: []
})

const selectedMethod = ref<'wechat' | 'alipay' | 'bank'>('wechat')
const paying = ref(false)

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    unpaid: '未支付',
    paid: '已支付',
    overdue: '已逾期',
    refunded: '已退款'
  }
  return map[status] || status
}

const isOverdue = (billInfo: Bill) => {
  if (billInfo.status === 'paid') return false
  return new Date(billInfo.dueDate) < new Date()
}

const selectMethod = (method: 'wechat' | 'alipay' | 'bank') => {
  selectedMethod.value = method
}

const handleBack = () => {
  uni.navigateBack()
}

const handlePay = async () => {
  if (!selectedMethod.value) {
    uni.showToast({
      title: '请选择支付方式',
      icon: 'none'
    })
    return
  }

  paying.value = true

  // 模拟支付流程
  setTimeout(() => {
    paying.value = false

    if (selectedMethod.value === 'wechat') {
      // 调用微信支付
      uni.requestPayment({
        provider: 'wxpay',
        timeStamp: String(Date.now()),
        nonceStr: 'random_string',
        package: 'prepay_id=xxx',
        signType: 'MD5',
        paySign: 'sign',
        success: () => {
          uni.showToast({
            title: '支付成功',
            icon: 'success'
          })
          setTimeout(() => {
            uni.navigateBack()
          }, 1500)
        },
        fail: () => {
          uni.showToast({
            title: '支付失败',
            icon: 'none'
          })
        }
      })
    } else if (selectedMethod.value === 'alipay') {
      // 调用支付宝支付
      uni.requestPayment({
        provider: 'alipay',
        orderInfo: 'order_info',
        success: () => {
          uni.showToast({
            title: '支付成功',
            icon: 'success'
          })
          setTimeout(() => {
            uni.navigateBack()
          }, 1500)
        },
        fail: () => {
          uni.showToast({
            title: '支付失败',
            icon: 'none'
          })
        }
      })
    } else {
      // 银行转账，显示转账信息
      uni.showModal({
        title: '银行转账信息',
        content: '开户行：中国银行\n账号：6217 0000 0000 0000 000\n户名：养老院有限公司\n\n转账后请备注账单号',
        confirmText: '我已转账',
        success: (res) => {
          if (res.confirm) {
            uni.showToast({
              title: '请等待财务确认',
              icon: 'none'
            })
          }
        }
      })
    }
  }, 500)
}

onLoad((options: any) => {
  const id = options.id
  // TODO: 根据ID加载账单详情
  // 模拟数据
  bill.value = {
    id: parseInt(id) || 1,
    billNo: 'B202603001',
    billType: '床位费',
    period: '2026年3月',
    billDate: '2026-03-01',
    dueDate: '2026-03-10',
    amount: 3500,
    status: 'unpaid',
    details: [
      { name: '床位费', description: '双人间 (3号楼201)', amount: 2800 },
      { name: '护理费', description: '二级护理', amount: 500 },
      { name: '伙食费', description: '标准餐', amount: 200 }
    ]
  }
})
</script>

<style lang="scss" scoped>
.payment-page {
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

.payment-content {
  padding: 20rpx 30rpx;
}

.bill-section,
.detail-section,
.payment-section {
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20rpx;
}

.bill-card,
.detail-list,
.payment-methods {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
}

.bill-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.bill-no {
  font-size: 26rpx;
  color: #606266;
}

.bill-status {
  padding: 8rpx 20rpx;
  border-radius: 8rpx;
  font-size: 24rpx;

  &.status-unpaid { background: #fdf6ec; color: #e6a23c; }
  &.status-paid { background: #f0f9ff; color: #67c23a; }
  &.status-overdue { background: #fef0f0; color: #f56c6c; }
}

.bill-amount {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  padding: 30rpx 0;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 20rpx;
}

.amount-label {
  font-size: 26rpx;
  color: #606266;
}

.amount-value {
  font-size: 48rpx;
  font-weight: 600;
  color: #f56c6c;
}

.bill-details {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  font-size: 26rpx;
}

.detail-label {
  color: #909399;
}

.detail-value {
  color: #303133;
}

.text-danger {
  color: #f56c6c;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20rpx 0;
  border-bottom: 1px solid #f5f7fa;

  &:last-child {
    border-bottom: none;
  }
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.item-name {
  font-size: 28rpx;
  color: #303133;
}

.item-desc {
  font-size: 24rpx;
  color: #909399;
}

.item-amount {
  font-size: 28rpx;
  font-weight: 500;
  color: #303133;
}

.detail-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 0;
  border-top: 1px solid #ebeef5;
  margin-top: 12rpx;
}

.total-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
}

.total-value {
  font-size: 36rpx;
  font-weight: 600;
  color: #f56c6c;
}

.method-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx;
  background: #f5f7fa;
  border-radius: 12rpx;
  margin-bottom: 16rpx;
  border: 2px solid transparent;
  transition: all 0.3s;

  &.active {
    border-color: #409eff;
    background: #ecf5ff;
  }

  &:last-child {
    margin-bottom: 0;
  }
}

.method-icon {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
}

.method-info {
  flex: 1;
}

.method-name {
  font-size: 28rpx;
  color: #303133;
}

.method-check {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.notice-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
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

.payment-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20rpx;
  padding: 20rpx 30rpx;
  background: #fff;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.payment-total {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
}

.total-label {
  font-size: 26rpx;
  color: #606266;
}

.total-amount {
  font-size: 40rpx;
  font-weight: 600;
  color: #f56c6c;
}

.pay-btn {
  flex: 1;
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
