<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card" @click="$router.push('/elderly')">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409eff">
              <el-icon size="30"><User /></el-icon>
            </div>
            <div class="stat-text">
              <div class="stat-value">{{ stats.elderly_count }}</div>
              <div class="stat-label">在院老人</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" @click="$router.push('/care')">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67c23a">
              <el-icon size="30"><Notebook /></el-icon>
            </div>
            <div class="stat-text">
              <div class="stat-value">{{ stats.today_care_count }}</div>
              <div class="stat-label">今日护理</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" @click="$router.push('/bills')">
          <div class="stat-content">
            <div class="stat-icon" style="background: #e6a23c">
              <el-icon size="30"><Wallet /></el-icon>
            </div>
            <div class="stat-text">
              <div class="stat-value">{{ stats.pending_bills }}</div>
              <div class="stat-label">待缴费账单</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" @click="$router.push('/staff')">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f56c6c">
              <el-icon size="30"><Avatar /></el-icon>
            </div>
            <div class="stat-text">
              <div class="stat-value">{{ stats.staff_count }}</div>
              <div class="stat-label">在职员工</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 入住率和服务请求 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <el-card class="occupancy-card">
          <div class="occupancy-content">
            <el-progress
              type="dashboard"
              :percentage="Math.round(stats.bed_occupancy_rate)"
              :color="getOccupancyColor(stats.bed_occupancy_rate)"
              :width="120"
            />
            <div class="occupancy-label">入住率</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="alert-card" @click="$router.push('/service')">
          <div class="alert-content">
            <el-badge :value="stats.pending_services" :max="99" class="alert-badge">
              <el-icon size="48" color="#f56c6c"><Bell /></el-icon>
            </el-badge>
            <div class="alert-label">待处理服务请求</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>财务概览</span>
          </template>
          <el-row :gutter="20" v-loading="financeLoading">
            <el-col :span="8">
              <div class="finance-item">
                <div class="finance-value">¥{{ formatMoney(financeStats.total_income) }}</div>
                <div class="finance-label">总收入</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="finance-item">
                <div class="finance-value success">¥{{ formatMoney(financeStats.monthly_income) }}</div>
                <div class="finance-label">本月收入</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="finance-item">
                <div class="finance-value warning">¥{{ formatMoney(financeStats.pending_amount) }}</div>
                <div class="finance-label">待收金额</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <!-- 待办和护理记录 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>待办事项</span>
              <el-button type="primary" link @click="$router.push('/service')">查看全部</el-button>
            </div>
          </template>
          <div v-if="stats.pending_tasks && stats.pending_tasks.length > 0">
            <div
              class="task-item"
              v-for="task in stats.pending_tasks"
              :key="task.id"
            >
              <div class="task-info">
                <el-tag :type="getTaskType(task.type)" size="small">{{ task.type }}</el-tag>
                <span class="task-elderly">{{ task.elderly_name }}</span>
              </div>
              <div class="task-time">{{ task.requested_at }}</div>
            </div>
          </div>
          <el-empty v-else description="暂无待办" :image-size="80" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最新护理记录</span>
              <el-button type="primary" link @click="$router.push('/care')">查看全部</el-button>
            </div>
          </template>
          <div v-if="stats.recent_cares && stats.recent_cares.length > 0">
            <div
              class="care-item"
              v-for="care in stats.recent_cares"
              :key="care.id"
            >
              <div class="care-info">
                <span class="care-name">{{ care.care_item }}</span>
                <span class="care-elderly">{{ care.elderly_name }}</span>
              </div>
              <div class="care-meta">
                <span class="care-staff">{{ care.staff_name }}</span>
                <span class="care-time">{{ care.recorded_at }}</span>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无记录" :image-size="80" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import instance from '@/api'

interface DashboardStats {
  elderly_count: number
  today_care_count: number
  pending_bills: number
  staff_count: number
  pending_services: number
  bed_occupancy_rate: number
  recent_cares: any[]
  pending_tasks: any[]
}

interface FinanceStats {
  total_income: number
  monthly_income: number
  pending_amount: number
}

const loading = ref(false)
const financeLoading = ref(false)
const stats = ref<DashboardStats>({
  elderly_count: 0,
  today_care_count: 0,
  pending_bills: 0,
  staff_count: 0,
  pending_services: 0,
  bed_occupancy_rate: 0,
  recent_cares: [],
  pending_tasks: []
})
const financeStats = ref<FinanceStats>({
  total_income: 0,
  monthly_income: 0,
  pending_amount: 0
})

const loadStats = async () => {
  loading.value = true
  try {
    const result = await instance.get('/stats/dashboard')
    stats.value = result || stats.value
  } catch (error) {
    console.error('加载统计失败', error)
  } finally {
    loading.value = false
  }
}

const loadFinanceStats = async () => {
  financeLoading.value = true
  try {
    const result = await instance.get('/stats/finance')
    financeStats.value = result || financeStats.value
  } catch (error) {
    console.error('加载财务统计失败', error)
  } finally {
    financeLoading.value = false
  }
}

const getOccupancyColor = (rate: number) => {
  if (rate >= 90) return '#67c23a'
  if (rate >= 70) return '#409eff'
  if (rate >= 50) return '#e6a23c'
  return '#f56c6c'
}

const formatMoney = (value: number) => {
  if (!value) return '0.00'
  return value.toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

const getTaskType = (type: string) => {
  const types: Record<string, string> = {
    '紧急': 'danger',
    '护理': 'primary',
    '医疗': 'warning',
    '送餐': 'success'
  }
  return types[type] || 'info'
}

onMounted(() => {
  loadStats()
  loadFinanceStats()
})
</script>

<style scoped>
.stat-card {
  cursor: pointer;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin-right: 20px;
}

.stat-text {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 5px;
}

.occupancy-card, .alert-card {
  height: 100%;
  cursor: pointer;
}

.occupancy-card:hover, .alert-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.occupancy-content, .alert-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 0;
}

.occupancy-label, .alert-label {
  margin-top: 10px;
  font-size: 14px;
  color: #666;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.finance-item {
  text-align: center;
  padding: 10px 0;
}

.finance-value {
  font-size: 20px;
  font-weight: bold;
  color: #333;
}

.finance-value.success {
  color: #67c23a;
}

.finance-value.warning {
  color: #e6a23c;
}

.finance-label {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.task-item, .care-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.task-item:last-child, .care-item:last-child {
  border-bottom: none;
}

.task-info, .care-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.task-elderly, .care-elderly {
  font-size: 14px;
  color: #333;
}

.task-time, .care-time {
  font-size: 12px;
  color: #999;
}

.care-info {
  flex-direction: column;
  align-items: flex-start;
}

.care-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.care-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.care-staff {
  font-size: 12px;
  color: #666;
}
</style>
