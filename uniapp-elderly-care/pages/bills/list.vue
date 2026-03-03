<template>
  <view class="bills-page">
    <!-- 页面头部 -->
    <view class="page-header">
      <text class="header-title">费用账单</text>
      <view class="header-actions">
        <view class="action-btn" @click="handleFilter">
          <uni-icons type="settings" size="14" />
          <text>筛选</text>
        </view>
      </view>
    </view>

    <!-- 欠费统计 -->
    <view class="bill-stats">
      <view class="stat-item" @click="handleStatClick('unpaid')">
        <text class="stat-value">¥3,200.00</text>
        <text class="stat-label">待支付</text>
      </view>
      <view class="stat-divider"></view>
      <view class="stat-item" @click="handleStatClick('paid')">
        <text class="stat-value stat-paid">¥8,500.00</text>
        <text class="stat-label">已支付</text>
      </view>
    </view>

    <!-- 账单列表 -->
    <scroll-view class="bills-list" scroll-y @scrolltolower="loadMore">
      <view
        class="bill-item"
        v-for="bill in billList"
        :key="bill.id"
        @click="handleBillDetail(bill)"
      >
        <!-- 头部 -->
        <view class="bill-header">
          <view class="bill-type-tag">{{ bill.billType }}</view>
          <view class="bill-status" :class="`status-${bill.status}`">
            {{ getStatusText(bill.status) }}
          </view>
        </view>

        <!-- 金额 -->
        <view class="bill-amount">
          <text class="amount-value">¥{{ bill.amount.toFixed(2) }}</text>
          <text class="amount-label">{{ bill.period }}</text>
        </view>

        <!-- 详情信息 -->
        <view class="bill-info">
          <view class="info-row">
            <text class="info-label">账单编号：</text>
            <text class="info-value">{{ bill.billNo }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">账单日期：</text>
            <text class="info-value">{{ bill.billDate }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">应付日期：</text>
            <text class="info-value" :class="{ 'text-danger': isOverdue(bill) }">
              {{ bill.dueDate }}
            </text>
          </view>
        </view>

        <!-- 底部操作 -->
        <view class="bill-footer" @click.stop>
          <button
            v-if="bill.status === 'unpaid'"
            class="pay-btn"
            @click="handlePay(bill)"
          >
            立即支付
          </button>
          <button v-else class="view-btn">
            查看详情
          </button>
        </view>
      </view>

      <!-- 加载更多 -->
      <view class="load-more" v-if="hasMore" @click="loadMore">
        <text>加载更多</text>
      </view>

      <!-- 空状态 -->
      <view v-if="!loading && billList.length === 0" class="empty-state">
        <text class="empty-text">暂无账单</text>
      </view>

      <!-- 加载中 -->
      <view v-if="loading" class="loading-state">
        <text>加载中...</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'

interface Bill {
  id: number
  billNo: string
  elderlyId: number
  elderlyName: string
  billType: string
  amount: number
  period: string
  billDate: string
  dueDate: string
  status: string
}

const billList = ref<Bill[]>([])
const loading = ref(false)
const hasMore = ref(false)
const pagination = reactive({
  page: 1,
  pageSize: 20
})

const currentFilter = ref('') // 'all', 'unpaid', 'paid'

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    unpaid: '未支付',
    paid: '已支付',
    overdue: '已逾期',
    refunded: '已退款'
  }
  return map[status] || status
}

const isOverdue = (bill: Bill) => {
  if (bill.status === 'paid') return false
  return new Date(bill.dueDate) < new Date()
}

const handleFilter = () => {
  uni.showActionSheet({
    itemList: ['全部', '未支付', '已支付', '已逾期'],
    success: (res) => {
      const filters = ['', 'unpaid', 'paid', 'overdue']
      currentFilter.value = filters[res.tapIndex || 0]
      loadBills()
    }
  })
}

const handleStatClick = (type: string) => {
  currentFilter.value = type === 'paid' ? 'paid' : 'unpaid'
  loadBills()
}

const handleBillDetail = (bill: Bill) => {
  uni.navigateTo({
    url: `/pages/bills/detail?id=${bill.id}`
  })
}

const handlePay = (bill: Bill) => {
  uni.navigateTo({
    url: `/pages/bills/payment?id=${bill.id}`
  })
}

const loadBills = async () => {
  loading.value = true
  try {
    // TODO: 实际 API 请求
    await new Promise(resolve => setTimeout(resolve, 500))

    // 模拟数据
    billList.value = [
      {
        id: 1,
        billNo: 'B202603001',
        elderlyId: 1,
        elderlyName: '张奶奶',
        billType: '床位费',
        amount: 3500,
        period: '2026年3月',
        billDate: '2026-03-01',
        dueDate: '2026-03-10',
        status: 'unpaid'
      },
      {
        id: 2,
        billNo: 'B202603002',
        elderlyId: 1,
        elderlyName: '张奶奶',
        billType: '护理费',
        amount: 1800,
        period: '2026年3月',
        billDate: '2026-03-01',
        dueDate: '2026-03-10',
        status: 'paid'
      },
      {
        id: 3,
        billNo: 'B202603003',
        elderlyId: 2,
        elderlyName: '王爷爷',
        billType: '床位费',
        amount: 3500,
        period: '2026年3月',
        billDate: '2026-03-01',
        dueDate: '2026-03-10',
        status: 'overdue'
      }
    ]

    // 根据筛选条件过滤
    if (currentFilter.value && currentFilter.value !== 'all') {
      billList.value = billList.value.filter(b => b.status === currentFilter.value)
    }

    hasMore.value = billList.value.length >= pagination.pageSize
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  if (!hasMore.value) return
  pagination.page++
  loadBills()
}

onMounted(() => {
  loadBills()
})
</script>

<style lang="scss" scoped>
.bills-page {
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

.header-actions {
  display: flex;
  gap: 20rpx;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 26rpx;
  color: #606266;
}

.bill-stats {
  display: flex;
  padding: 30rpx;
  background: #fff;
  margin-bottom: 20rpx;
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.stat-value {
  font-size: 40rpx;
  font-weight: 600;
  color: #f56c6c;

  &.stat-paid {
    color: #67c23a;
  }
}

.stat-label {
  font-size: 24rpx;
  color: #909399;
}

.stat-divider {
  width: 1rpx;
  height: 60rpx;
  background: #ebeef5;
}

.bills-list {
  flex: 1;
  padding: 0 30rpx 30rpx;
}

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
  margin-bottom: 20rpx;
}

.bill-type-tag {
  padding: 8rpx 20rpx;
  background: #ecf5ff;
  color: #409eff;
  font-size: 24rpx;
  border-radius: 8rpx;
}

.bill-status {
  padding: 8rpx 20rpx;
  border-radius: 8rpx;
  font-size: 24rpx;

  &.status-unpaid {
    background: #fdf6ec;
    color: #e6a23c;
  }

  &.status-paid {
    background: #f0f9ff;
    color: #67c23a;
  }

  &.status-overdue {
    background: #fef0f0;
    color: #f56c6c;
  }

  &.status-refunded {
    background: #f4f4f5;
    color: #909399;
  }
}

.bill-amount {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
  margin-bottom: 20rpx;
}

.amount-value {
  font-size: 44rpx;
  font-weight: 600;
  color: #f56c6c;
}

.amount-label {
  font-size: 24rpx;
  color: #909399;
}

.bill-info {
  margin-bottom: 20rpx;
}

.info-row {
  display: flex;
  align-items: center;
  font-size: 26rpx;
  line-height: 1.8;

  &:last-child {
    margin-bottom: 0;
  }
}

.info-label {
  color: #909399;
  margin-right: 8rpx;
}

.info-value {
  color: #303133;
  flex: 1;
}

.text-danger {
  color: #f56c6c;
}

.bill-footer {
  text-align: right;
}

.pay-btn {
  padding: 16rpx 48rpx;
  background: #409eff;
  color: #fff;
  font-size: 28rpx;
  border: none;
  border-radius: 8rpx;
}

.view-btn {
  padding: 16rpx 48rpx;
  background: #f5f7fa;
  color: #606266;
  font-size: 28rpx;
  border: none;
  border-radius: 8rpx;
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
