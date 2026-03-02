<template>
  <div class="bills-page">
    <van-nav-bar title="费用账单" />

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list v-model:loading="loading" :finished="finished" @load="onLoad">
        <div class="bill-list">
          <div class="bill-item" v-for="item in list" :key="item.id" @click="goToDetail(item.id)">
            <div class="bill-header">
              <span class="bill-no">{{ item.bill_no }}</span>
              <van-tag :type="item.status === 'paid' ? 'success' : 'warning'">
                {{ item.status === 'paid' ? '已支付' : '待支付' }}
              </van-tag>
            </div>
            <div class="bill-body">
              <div class="amount">¥{{ item.total_amount }}</div>
              <div class="period">{{ formatPeriod(item.bill_period_start, item.bill_period_end) }}</div>
            </div>
            <div class="bill-footer">
              <span class="time">{{ formatTime(item.created_at) }}</span>
              <van-icon name="arrow" />
            </div>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <van-empty v-if="list.length === 0 && !loading" description="暂无账单" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'

const router = useRouter()
const list = ref<any[]>([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)

onMounted(() => {
  // 默认加载第一个老人的账单
  loadBills()
})

function loadBills() {
  loading.value = true
  api.bill.list({ page: 1, page_size: 20 }).then((data: any) => {
    list.value = data.list || []
    loading.value = false
    finished.value = true
  }).catch(() => {
    loading.value = false
  })
}

function onLoad() {
  loadBills()
}

function onRefresh() {
  finished.value = false
  loadBills()
  refreshing.value = false
}

function goToDetail(id: number) {
  router.push(`/bills/${id}`)
}

function formatPeriod(start: string, end: string) {
  if (!start || !end) return ''
  const s = new Date(start)
  const e = new Date(end)
  return `${s.getMonth() + 1}月${s.getDate()}日 - ${e.getMonth() + 1}月${e.getDate()}日`
}

function formatTime(time: string) {
  const date = new Date(time)
  return `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()}`
}
</script>

<style scoped>
.bills-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.bill-list {
  padding: 15px;
}

.bill-item {
  background: #fff;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 12px;
}

.bill-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.bill-no {
  font-size: 14px;
  color: #999;
}

.bill-body {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.amount {
  font-size: 24px;
  font-weight: 600;
  color: #f56c6c;
}

.period {
  font-size: 13px;
  color: #999;
}

.bill-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #eee;
}

.bill-footer .time {
  font-size: 12px;
  color: #999;
}
</style>
