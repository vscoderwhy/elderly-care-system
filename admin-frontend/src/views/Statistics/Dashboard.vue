<template>
  <el-row :gutter="20">
    <!-- 统计卡片 -->
    <el-col :span="6" v-for="stat in stats" :key="stat.key">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon" :style="{ backgroundColor: stat.color }">
            <component :is="stat.icon" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </el-card>
    </el-col>
  </el-row>

  <el-row :gutter="20" style="margin-top: 20px">
    <!-- 入住率趋势 -->
    <el-col :span="12">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>入住率趋势</span>
            <el-select v-model="occupancyDays" size="small" style="width: 100px" @change="fetchOccupancyTrend">
              <el-option label="7天" :value="7" />
              <el-option label="30天" :value="30" />
              <el-option label="90天" :value="90" />
            </el-select>
          </div>
        </template>
        <div ref="occupancyChart" style="height: 300px"></div>
      </el-card>
    </el-col>

    <!-- 护理等级分布 -->
    <el-col :span="12">
      <el-card>
        <template #header>
          <span>护理等级分布</span>
        </template>
        <div ref="careLevelChart" style="height: 300px"></div>
      </el-card>
    </el-col>
  </el-row>

  <el-row :gutter="20" style="margin-top: 20px">
    <!-- 性别分布 -->
    <el-col :span="8">
      <el-card>
        <template #header>
          <span>性别分布</span>
        </template>
        <div ref="genderChart" style="height: 250px"></div>
      </el-card>
    </el-col>

    <!-- 年龄分布 -->
    <el-col :span="8">
      <el-card>
        <template #header>
          <span>年龄分布</span>
        </template>
        <div ref="ageChart" style="height: 250px"></div>
      </el-card>
    </el-col>

    <!-- 今日护理统计 -->
    <el-col :span="8">
      <el-card>
        <template #header>
          <span>今日护理统计</span>
        </template>
        <div ref="careStatsChart" style="height: 250px"></div>
      </el-card>
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, markRaw } from 'vue'
import axios from '@/api/index'
import * as echarts from 'echarts'
import { User, Bed, Calendar, TrendCharts } from '@element-plus/icons-vue'

const occupancyChart = ref<HTMLElement>()
const careLevelChart = ref<HTMLElement>()
const genderChart = ref<HTMLElement>()
const ageChart = ref<HTMLElement>()
const careStatsChart = ref<HTMLElement>()

let occupancyChartInstance: echarts.ECharts | null = null
let careLevelChartInstance: echarts.ECharts | null = null
let genderChartInstance: echarts.ECharts | null = null
let ageChartInstance: echarts.ECharts | null = null
let careStatsChartInstance: echarts.ECharts | null = null

const occupancyDays = ref(30)

const stats = reactive([
  { key: 'elderly', label: '在院老人', value: 0, icon: markRaw(User), color: '#409eff' },
  { key: 'beds', label: '总床位数', value: 0, icon: markRaw(Bed), color: '#67c23a' },
  { key: 'occupancy', label: '入住率', value: '0%', icon: markRaw(TrendCharts), color: '#e6a23c' },
  { key: 'care', label: '今日护理', value: 0, icon: markRaw(Calendar), color: '#f56c6c' }
])

const fetchDashboardStats = async () => {
  try {
    const res = await axios.get('/statistics/dashboard')
    const data = res.data

    stats[0].value = data.elderly_total || 0
    stats[1].value = data.bed_total || 0
    stats[2].value = ((data.occupancy_rate || 0).toFixed(1)) + '%'
    stats[3].value = data.care_records_today || 0

    // 渲染护理等级分布图
    renderCareLevelChart(data.care_level_dist || {})
    // 渲染性别分布图
    renderGenderChart(data.gender_dist || {})
  } catch (error) {
    console.error('Failed to fetch dashboard stats:', error)
  }
}

const fetchOccupancyTrend = async () => {
  try {
    const res = await axios.get('/statistics/occupancy-trend', {
      params: { days: occupancyDays.value }
    })
    const data = res.data
    renderOccupancyChart(data)
  } catch (error) {
    console.error('Failed to fetch occupancy trend:', error)
  }
}

const fetchAgeDistribution = async () => {
  try {
    const res = await axios.get('/statistics/age-distribution')
    const data = res.data
    renderAgeChart(data)
  } catch (error) {
    console.error('Failed to fetch age distribution:', error)
  }
}

const fetchCareStats = async () => {
  try {
    const res = await axios.get('/statistics/care', {
      params: {
        start_date: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        end_date: new Date().toISOString().split('T')[0]
      }
    })
    const data = res.data
    renderCareStatsChart(data.item_stats || {})
  } catch (error) {
    console.error('Failed to fetch care stats:', error)
  }
}

const renderOccupancyChart = (data: any[]) => {
  if (!occupancyChart.value) return

  if (!occupancyChartInstance) {
    occupancyChartInstance = echarts.init(occupancyChart.value)
  }

  const dates = data.map(d => d.date)
  const occupied = data.map(d => d.occupied)
  const total = data.map(d => d.total)

  const option = {
    tooltip: { trigger: 'axis' },
    legend: { data: ['已入住', '总床位'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: dates },
    yAxis: { type: 'value' },
    series: [
      { name: '已入住', type: 'line', data: occupied, smooth: true, itemStyle: { color: '#409eff' } },
      { name: '总床位', type: 'line', data: total, smooth: true, itemStyle: { color: '#67c23a' } }
    ]
  }

  occupancyChartInstance.setOption(option)
}

const renderCareLevelChart = (data: Record<string, number>) => {
  if (!careLevelChart.value) return

  if (!careLevelChartInstance) {
    careLevelChartInstance = echarts.init(careLevelChart.value)
  }

  const levels = ['1级护理', '2级护理', '3级护理', '4级护理', '5级护理', '特级护理']
  const values = levels.map(l => data[parseInt(l)] || 0)

  const option = {
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', left: 'left' },
    series: [{
      type: 'pie',
      radius: '60%',
      data: levels.map((l, i) => ({ name: l, value: values[i] })),
      emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.5)' } }
    }]
  }

  careLevelChartInstance.setOption(option)
}

const renderGenderChart = (data: Record<string, number>) => {
  if (!genderChart.value) return

  if (!genderChartInstance) {
    genderChartInstance = echarts.init(genderChart.value)
  }

  const option = {
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', left: 'left' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: [
        { name: '男', value: data.male || data['男'] || 0, itemStyle: { color: '#409eff' } },
        { name: '女', value: data.female || data['女'] || 0, itemStyle: { color: '#f56c6c' } }
      ]
    }]
  }

  genderChartInstance.setOption(option)
}

const renderAgeChart = (data: Record<string, number>) => {
  if (!ageChart.value) return

  if (!ageChartInstance) {
    ageChartInstance = echarts.init(ageChart.value)
  }

  const ages = Object.keys(data)
  const values = Object.values(data)

  const option = {
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ages },
    yAxis: { type: 'value' },
    series: [{
      type: 'bar',
      data: values,
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: '#83bff6' },
          { offset: 1, color: '#188df0' }
        ])
      }
    }]
  }

  ageChartInstance.setOption(option)
}

const renderCareStatsChart = (data: Record<string, number>) => {
  if (!careStatsChart.value) return

  if (!careStatsChartInstance) {
    careStatsChartInstance = echarts.init(careStatsChart.value)
  }

  // 取前5项
  const items = Object.keys(data).slice(0, 5)
  const values = items.map(k => data[k])

  const option = {
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: items },
    yAxis: { type: 'value' },
    series: [{
      type: 'bar',
      data: values,
      itemStyle: { color: '#67c23a' }
    }]
  }

  careStatsChartInstance.setOption(option)
}

const handleResize = () => {
  occupancyChartInstance?.resize()
  careLevelChartInstance?.resize()
  genderChartInstance?.resize()
  ageChartInstance?.resize()
  careStatsChartInstance?.resize()
}

onMounted(() => {
  fetchDashboardStats()
  fetchOccupancyTrend()
  fetchAgeDistribution()
  fetchCareStats()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  occupancyChartInstance?.dispose()
  careLevelChartInstance?.dispose()
  genderChartInstance?.dispose()
  ageChartInstance?.dispose()
  careStatsChartInstance?.dispose()
})
</script>

<style scoped>
.stat-card {
  margin-bottom: 20px;
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
  margin-right: 16px;
}

.stat-icon :deep(svg) {
  width: 30px;
  height: 30px;
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
