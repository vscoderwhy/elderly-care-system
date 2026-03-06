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
import { ref, reactive, onMounted, onUnmounted, markRaw, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import axios from '@/api/index'
import * as echarts from 'echarts'
import { User, Grid, Calendar, TrendCharts } from '@element-plus/icons-vue'

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
  { key: 'beds', label: '总床位数', value: 0, icon: markRaw(Grid), color: '#67c23a' },
  { key: 'occupancy', label: '入住率', value: '0%', icon: markRaw(TrendCharts), color: '#e6a23c' },
  { key: 'care', label: '今日护理', value: 0, icon: markRaw(Calendar), color: '#f56c6c' }
])

const fetchDashboardStats = async () => {
  try {
    const res = await axios.get('/statistics/dashboard')
    const data = res.data?.data || res.data || {}

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
    const data = res.data?.data || res.data || []
    renderOccupancyChart(Array.isArray(data) ? data : [])
  } catch (error) {
    console.error('Failed to fetch occupancy trend:', error)
  }
}

const fetchAgeDistribution = async () => {
  try {
    const res = await axios.get('/statistics/age-distribution')
    const data = res.data?.data || res.data || {}
    renderAgeChart(data && typeof data === 'object' ? data : {})
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
    const data = res.data?.data || res.data || {}
    renderCareStatsChart(data.item_stats || data.items || {})
  } catch (error) {
    console.error('Failed to fetch care stats:', error)
  }
}

const renderOccupancyChart = (data: any[]) => {
  if (!occupancyChart.value) {
    console.warn('occupancyChart container not ready')
    return
  }

  try {
    if (!occupancyChartInstance) {
      occupancyChartInstance = echarts.init(occupancyChart.value)
    }

    // 验证数据有效性
    if (!Array.isArray(data) || data.length === 0) {
      console.warn('Invalid occupancy trend data:', data)
      // 显示空状态提示
      occupancyChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        xAxis: { show: false },
        yAxis: { show: false },
        series: []
      })
      return
    }

    const dates = data?.map?.((d: any) => d?.date) || []
    const occupied = data?.map?.((d: any) => d?.occupied) || []
    const total = data?.map?.((d: any) => d?.total) || []

    // 验证数据不为空
    if (dates.length === 0 || occupied.length === 0 || total.length === 0) {
      console.warn('Empty occupancy trend arrays')
      occupancyChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        xAxis: { show: false },
        yAxis: { show: false },
        series: []
      })
      return
    }

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
  } catch (error) {
    console.error('Failed to render occupancy chart:', error)
    // 显示错误提示
    occupancyChartInstance?.setOption({
      title: { text: '加载失败', left: 'center', top: 'center', textStyle: { color: '#999' } },
      xAxis: { show: false },
      yAxis: { show: false },
      series: []
    })
  }
}

const renderCareLevelChart = (data: Record<string, number>) => {
  if (!careLevelChart.value) return

  try {
    if (!careLevelChartInstance) {
      careLevelChartInstance = echarts.init(careLevelChart.value)
    }

    // 验证数据有效性
    if (!data || typeof data !== 'object') {
      console.warn('Invalid care level data:', data)
      careLevelChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        series: []
      })
      return
    }

    const levels = ['1级护理', '2级护理', '3级护理', '4级护理', '5级护理', '特级护理']
    const values = levels.map(l => {
      // 尝试多种可能的键名格式
      const num = parseInt(l.replace('级护理', ''))
      return data[num] || data[l] || data[`level_${num}`] || 0
    })

    // 验证数据不为全0
    const total = values.reduce((a, b) => a + b, 0)
    if (total === 0) {
      console.warn('All care level values are zero')
      careLevelChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        series: []
      })
      return
    }

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
  } catch (error) {
    console.error('Failed to render care level chart:', error)
    careLevelChartInstance?.setOption({
      title: { text: '加载失败', left: 'center', top: 'center', textStyle: { color: '#999' } },
      series: []
    })
  }
}

const renderGenderChart = (data: Record<string, number>) => {
  if (!genderChart.value) return

  try {
    if (!genderChartInstance) {
      genderChartInstance = echarts.init(genderChart.value)
    }

    // 验证数据有效性
    if (!data || typeof data !== 'object') {
      console.warn('Invalid gender distribution data:', data)
      genderChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        series: []
      })
      return
    }

    // 支持多种可能的键名格式
    const male = data.male || data['男'] || data.Male || 0
    const female = data.female || data['女'] || data.Female || 0

    // 验证数据不为全0
    if (male === 0 && female === 0) {
      console.warn('Gender distribution values are all zero')
      genderChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        series: []
      })
      return
    }

    const option = {
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { name: '男', value: male, itemStyle: { color: '#409eff' } },
          { name: '女', value: female, itemStyle: { color: '#f56c6c' } }
        ]
      }]
    }

    genderChartInstance.setOption(option)
  } catch (error) {
    console.error('Failed to render gender chart:', error)
    genderChartInstance?.setOption({
      title: { text: '加载失败', left: 'center', top: 'center', textStyle: { color: '#999' } },
      series: []
    })
  }
}

const renderAgeChart = (data: Record<string, number>) => {
  if (!ageChart.value) return

  try {
    if (!ageChartInstance) {
      ageChartInstance = echarts.init(ageChart.value)
    }

    // 验证数据有效性
    if (!data || typeof data !== 'object') {
      console.warn('Invalid age distribution data:', data)
      ageChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        xAxis: { show: false },
        yAxis: { show: false },
        series: []
      })
      return
    }

    const ages = Object.keys(data)
    const values = Object.values(data)

    // 验证数据不为空
    if (!ages.length || !values.length || values.every(v => v === 0)) {
      console.warn('Empty age distribution data')
      ageChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        xAxis: { show: false },
        yAxis: { show: false },
        series: []
      })
      return
    }

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
  } catch (error) {
    console.error('Failed to render age chart:', error)
    ageChartInstance?.setOption({
      title: { text: '加载失败', left: 'center', top: 'center', textStyle: { color: '#999' } },
      xAxis: { show: false },
      yAxis: { show: false },
      series: []
    })
  }
}

const renderCareStatsChart = (data: Record<string, number>) => {
  if (!careStatsChart.value) return

  try {
    if (!careStatsChartInstance) {
      careStatsChartInstance = echarts.init(careStatsChart.value)
    }

    // 验证数据有效性
    if (!data || typeof data !== 'object') {
      console.warn('Invalid care stats data:', data)
      careStatsChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        xAxis: { show: false },
        yAxis: { show: false },
        series: []
      })
      return
    }

    // 取前5项（按数量降序）
    const items = Object.keys(data)
      .sort((a, b) => (data[b] || 0) - (data[a] || 0))
      .slice(0, 5)

    const values = items.map(k => data[k] || 0)

    // 验证数据不为空
    if (!items.length || values.every(v => v === 0)) {
      console.warn('Empty care stats data')
      careStatsChartInstance.setOption({
        title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#999' } },
        xAxis: { show: false },
        yAxis: { show: false },
        series: []
      })
      return
    }

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
  } catch (error) {
    console.error('Failed to render care stats chart:', error)
    careStatsChartInstance?.setOption({
      title: { text: '加载失败', left: 'center', top: 'center', textStyle: { color: '#999' } },
      xAxis: { show: false },
      yAxis: { show: false },
      series: []
    })
  }
}

const handleResize = () => {
  try {
    occupancyChartInstance?.resize()
    careLevelChartInstance?.resize()
    genderChartInstance?.resize()
    ageChartInstance?.resize()
    careStatsChartInstance?.resize()
  } catch (error) {
    console.error('Failed to resize charts:', error)
  }
}

onMounted(async () => {
  await nextTick()

  try {
    await Promise.all([
      fetchDashboardStats(),
      fetchOccupancyTrend(),
      fetchAgeDistribution(),
      fetchCareStats()
    ])
  } catch (error) {
    ElMessage.error('数据加载失败，请刷新重试')
    console.error('Failed to load data:', error)
  }

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  try {
    occupancyChartInstance?.dispose()
    careLevelChartInstance?.dispose()
    genderChartInstance?.dispose()
    ageChartInstance?.dispose()
    careStatsChartInstance?.dispose()
  } catch (error) {
    console.error('Failed to dispose charts:', error)
  }
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
